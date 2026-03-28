package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
)

type Cluesheet struct {
	Id         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Origin_id  *uuid.UUID `json:"origin_id"`  /* originating cluesheet */
	Created_by string     `json:"created_by"` /* ipa unique id */
	Created_at time.Time  `json:"created_at"` /* no timezone by default */
	Edited_by  string     `json:"edited_by"`  /* ipa unique id */
	Edited_at  time.Time  `json:"edited_at"`  /* no timezone */
	Visibility string     `json:"visibility"` /* TODO: int key? or defined enum? I think SQL enums are difficult to work with in migrations */
	Owners     []string   `json:"owners"`
	Groups     []string   `json:"groups"`
	Clues      *[]Clue    `json:"clues,omitzero"`
}

// RuleInfo is the resolved rule returned to clients, mapping clue_rules fields to frontend expectations.
type RuleInfo struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

type Clue struct {
	Id          uuid.UUID        `json:"id"`
	Description string           `json:"description"`
	Rule_id     *uuid.UUID       `json:"rule_id"`
	Rule_params *json.RawMessage `json:"rule_params"`
	Origin_id   uuid.UUID        `json:"origin_id"` /* originating cluesheet */
	Created_by  string           `json:"created_by"`
	Created_at  time.Time        `json:"created_at"`
	Edited_by   string           `json:"edited_by"`
	Edited_at   time.Time        `json:"edited_at"`
	Tags        []string         `json:"tags"`
	Children    []*Clue          `json:"children,omitzero"`
	Rule        *RuleInfo        `json:"rule,omitempty"`
	Completions *int             `json:"completions,omitempty"`
}

type UserParticipation struct {
	Cluesheet_id uuid.UUID `json:"cluesheet_id"`
	Ipa_uid      string    `json:"ipa_uid"`
	Hidden       bool      `json:"hidden"`
}

type UserProgress struct {
	Ipa_uid     string    `json:"ipa_uid"`
	Clue_id     uuid.UUID `json:"clue_id"`
	Completions int       `json:"completions"` // TODO this should become a double
}

type ClueRelation struct {
	Id        uuid.UUID `json:"id"`
	Parent_id uuid.UUID `json:"parent_id"`
	Child_id  uuid.UUID `json:"child_id"`
}

// clueRule is an internal scan target for clue_rules rows.
type clueRule struct {
	Id      uuid.UUID
	Name    string
	Jq_expr string
}

func GetClues(ctx context.Context, conn *pgxpool.Pool, cluesheet uuid.UUID) ([]Clue, error) {
	rows, err := conn.Query(ctx, `select * from clue where origin_id = $1`, cluesheet)
	if err != nil {
		return nil, fmt.Errorf("failed to query clues: %w", err)
	}
	clues, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Clue])

	var clueMap map[uuid.UUID]*Clue = make(map[uuid.UUID]*Clue, len(clues))
	var out map[uuid.UUID]*Clue = make(map[uuid.UUID]*Clue)

	for _, clue := range clues {
		clueMap[clue.Id] = &clue
		out[clue.Id] = &clue
	}
	var mex sync.Mutex
	var relations []ClueRelation
	b := &pgx.Batch{}

	for _, child := range clues {
		qq := b.Queue(`select parent_id, child_id from clue_relation where child_id = $1`, child.Id)
		qq.Query(func(rows pgx.Rows) error {
			rel, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[ClueRelation])
			// Empty rows are fine
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			} else if err != nil {
				return err
			}

			mex.Lock()
			defer mex.Unlock()
			relations = append(relations, rel)
			return nil
		})
	}
	br := conn.SendBatch(ctx, b)
	err = br.Close()
	if err != nil {
		return nil, fmt.Errorf("failed querying for clue relations: %w", err)
	}

	for _, rel := range relations {
		clueMap[rel.Parent_id].Children = append(clueMap[rel.Parent_id].Children, clueMap[rel.Child_id])
		delete(out, rel.Child_id)
	}

	// Resolve rules for clues that have a rule_id
	rb := &pgx.Batch{}
	type ruleResult struct {
		clueID uuid.UUID
		rule   clueRule
	}
	var ruleResults []ruleResult
	var ruleMex sync.Mutex

	for id, clue := range clueMap {
		if clue.Rule_id == nil {
			continue
		}
		clueID := id
		ruleID := *clue.Rule_id
		qq := rb.Queue(`select id, name, jq_expr from clue_rules where id = $1`, ruleID)
		qq.Query(func(rows pgx.Rows) error {
			rule, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[clueRule])
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			} else if err != nil {
				return err
			}
			ruleMex.Lock()
			defer ruleMex.Unlock()
			ruleResults = append(ruleResults, ruleResult{clueID: clueID, rule: rule})
			return nil
		})
	}
	rbr := conn.SendBatch(ctx, rb)
	err = rbr.Close()
	if err != nil {
		return nil, fmt.Errorf("failed querying for clue rules: %w", err)
	}

	for _, rr := range ruleResults {
		clueMap[rr.clueID].Rule = &RuleInfo{Key: rr.rule.Name, Description: rr.rule.Jq_expr}
	}

	return slices.Collect(PointerToValue[iter.Seq[*Clue], Clue](maps.Values(out))), nil
}

// GetCluesForUser returns the clue tree for a cluesheet with each clue's completions for the given user populated.
func GetCluesForUser(ctx context.Context, conn *pgxpool.Pool, cluesheet uuid.UUID, ipaUID string) ([]Clue, error) {
	clues, err := GetClues(ctx, conn, cluesheet)
	if err != nil {
		return nil, err
	}

	// Collect all clue IDs (including children at all depths)
	var allClues []*Clue
	var collectAll func(clues []Clue)
	collectAll = func(clues []Clue) {
		for i := range clues {
			allClues = append(allClues, &clues[i])
			var children []Clue
			for _, c := range clues[i].Children {
				children = append(children, *c)
			}
			collectAll(children)
		}
	}
	collectAll(clues)

	// Batch-query user_progress for all clue IDs
	b := &pgx.Batch{}
	type progressResult struct {
		clue        *Clue
		completions int
	}
	var progressResults []progressResult
	var mex sync.Mutex

	for _, clue := range allClues {
		c := clue
		qq := b.Queue(`select ipa_uid, clue_id, completions from user_progress where clue_id = $1 and ipa_uid = $2`, c.Id, ipaUID)
		qq.Query(func(rows pgx.Rows) error {
			progress, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[UserProgress])
			if errors.Is(err, pgx.ErrNoRows) {
				mex.Lock()
				defer mex.Unlock()
				progressResults = append(progressResults, progressResult{clue: c, completions: 0})
				return nil
			} else if err != nil {
				return err
			}
			mex.Lock()
			defer mex.Unlock()
			progressResults = append(progressResults, progressResult{clue: c, completions: progress.Completions})
			return nil
		})
	}
	br := b.Queue("") // flush
	_ = br
	batchResults := conn.SendBatch(ctx, b)
	err = batchResults.Close()
	if err != nil {
		return nil, fmt.Errorf("failed querying user progress: %w", err)
	}

	for _, pr := range progressResults {
		completions := pr.completions
		pr.clue.Completions = &completions
	}

	return clues, nil
}

func PointerToValue[Iter iter.Seq[*V], V any](in Iter) iter.Seq[V] {
	return func(yeild func(V) bool) {
		for v := range in {
			if !yeild(*v) {
				return
			}
		}
	}
}
