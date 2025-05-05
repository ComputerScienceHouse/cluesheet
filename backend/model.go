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
	Id         uuid.UUID
	Name       string
	Origin_id  *uuid.UUID /* originating cluesheet */
	Created_by string     /* ipa unique id */
	Created_at time.Time  /* no timezone by default */
	Edited_by  string     /* ipa unique id */
	Edited_at  time.Time  /* no timezone */
	Visibility string     /* TODO: int key? or defined enum? I think SQL enums are difficult to work with in migrations */
	Owners     []string
	Groups     []string
	Clues      *[]Clue
}

type Clue struct {
	Id          uuid.UUID
	Description string
	Rule_id     *uuid.UUID
	Rule_params *json.RawMessage
	Origin_id   uuid.UUID /* originating cluesheet */
	Created_by  string
	Created_at  time.Time
	Edited_by   string
	Edited_at   time.Time
	Tags        []string
	Children    []*Clue
}

type ClueRelation struct {
	Id        uuid.UUID
	Parent_id uuid.UUID
	Child_id  uuid.UUID
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

	return slices.Collect(PointerToValue[iter.Seq[*Clue], Clue](maps.Values(out))), nil
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
