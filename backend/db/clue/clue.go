package dbclue

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"sync"

	"csh/cluesheet/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// clueRule is an internal scan target for clue_rules rows.
type clueRule struct {
	Id      uuid.UUID
	Name    string
	Jq_expr string
}

func pointerToValue[Iter iter.Seq[*V], V any](in Iter) iter.Seq[V] {
	return func(yeild func(V) bool) {
		for v := range in {
			if !yeild(*v) {
				return
			}
		}
	}
}

// CreateClueResult holds the result of a CreateClue call.
type CreateClueResult struct {
	Clue         model.Clue
	ClueRelation *model.ClueRelation `json:"ClueRelation,omitempty"`
}

func GetClues(ctx context.Context, conn *pgxpool.Pool, cluesheetID uuid.UUID) ([]model.Clue, error) {
	rows, err := conn.Query(ctx, `select * from clue where origin_id = $1`, cluesheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query clues: %w", err)
	}
	clues, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Clue])

	var clueMap map[uuid.UUID]*model.Clue = make(map[uuid.UUID]*model.Clue, len(clues))
	var out map[uuid.UUID]*model.Clue = make(map[uuid.UUID]*model.Clue)

	for _, clue := range clues {
		clueMap[clue.Id] = &clue
		out[clue.Id] = &clue
	}
	var mex sync.Mutex
	var relations []model.ClueRelation
	b := &pgx.Batch{}

	for _, child := range clues {
		qq := b.Queue(`select parent_id, child_id from clue_relation where child_id = $1`, child.Id)
		qq.Query(func(rows pgx.Rows) error {
			rel, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[model.ClueRelation])
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
		clueMap[rr.clueID].Rule = &model.RuleInfo{Key: rr.rule.Name, Description: rr.rule.Jq_expr}
	}

	return slices.Collect(pointerToValue[iter.Seq[*model.Clue], model.Clue](maps.Values(out))), nil
}

// GetCluesForUser returns the clue tree for a cluesheet with each clue's completions for the given user populated.
func GetCluesForUser(ctx context.Context, conn *pgxpool.Pool, cluesheetID uuid.UUID, ipaUID string) ([]model.Clue, error) {
	clues, err := GetClues(ctx, conn, cluesheetID)
	if err != nil {
		return nil, err
	}

	// Collect all clue IDs (including children at all depths)
	var allClues []*model.Clue
	var collectAll func(clues []model.Clue)
	collectAll = func(clues []model.Clue) {
		for i := range clues {
			allClues = append(allClues, &clues[i])
			var children []model.Clue
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
		clue        *model.Clue
		completions int
	}
	var progressResults []progressResult
	var mex sync.Mutex

	for _, clue := range allClues {
		c := clue
		qq := b.Queue(`select ipa_uid, clue_id, completions from user_progress where clue_id = $1 and ipa_uid = $2`, c.Id, ipaUID)
		qq.Query(func(rows pgx.Rows) error {
			progress, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[model.UserProgress])
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

func CreateClue(ctx context.Context, conn *pgxpool.Pool, clue model.Clue, parentClueID *uuid.UUID) (CreateClueResult, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return CreateClueResult{}, fmt.Errorf("create clue begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `insert into clue(id, description, tags, origin_id, created_by, created_at, edited_by, edited_at) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		clue.Id,
		clue.Description,
		clue.Tags,
		clue.Origin_id,
		clue.Created_by,
		clue.Created_at,
		clue.Edited_by,
		clue.Edited_at,
	)
	if err != nil {
		return CreateClueResult{}, fmt.Errorf("create clue insert: %w", err)
	}

	var cr *model.ClueRelation

	if parentClueID != nil {
		cr = &model.ClueRelation{
			Id:        uuid.New(),
			Parent_id: *parentClueID,
			Child_id:  clue.Id,
		}

		_, err = tx.Exec(ctx, `insert into clue_relation(id, parent_id, child_id) values ($1, $2, $3)`, cr.Id, cr.Parent_id, cr.Child_id)
		if err != nil {
			return CreateClueResult{}, fmt.Errorf("create clue relation insert: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CreateClueResult{}, fmt.Errorf("create clue commit: %w", err)
	}

	return CreateClueResult{Clue: clue, ClueRelation: cr}, nil
}
