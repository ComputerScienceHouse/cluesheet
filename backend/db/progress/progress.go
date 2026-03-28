package dbprogress

import (
	"context"
	"errors"
	"fmt"

	"csh/cluesheet/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetProgress(ctx context.Context, conn *pgxpool.Pool, clueID uuid.UUID, ipaUID string) (model.UserProgress, error) {
	rows, err := conn.Query(ctx, `select * from user_progress where clue_id = $1 and ipa_uid = $2`, clueID, ipaUID)
	if err != nil {
		return model.UserProgress{}, fmt.Errorf("get progress: %w", err)
	}

	progress, err := pgx.CollectOneRow[model.UserProgress](rows, pgx.RowToStructByNameLax[model.UserProgress])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.UserProgress{Ipa_uid: ipaUID, Clue_id: clueID, Completions: 0}, nil
	} else if err != nil {
		return model.UserProgress{}, fmt.Errorf("get progress: %w", err)
	}

	return progress, nil
}

func UpsertProgress(ctx context.Context, conn *pgxpool.Pool, progress model.UserProgress) error {
	_, err := conn.Exec(ctx, `insert into user_progress (ipa_uid, clue_id, completions) values ($1, $2, $3) on conflict (ipa_uid, clue_id) do update set completions = $3`, progress.Ipa_uid, progress.Clue_id, progress.Completions)
	if err != nil {
		return fmt.Errorf("upsert progress: %w", err)
	}

	return nil
}
