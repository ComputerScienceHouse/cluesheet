package dbparticipation

import (
	"context"
	"errors"
	"fmt"

	"csh/cluesheet/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetParticipation(ctx context.Context, conn *pgxpool.Pool, cluesheetID uuid.UUID, ipaUID string) (model.UserParticipation, error) {
	rows, err := conn.Query(ctx, `select * from user_participation where cluesheet_id = $1 and ipa_uid = $2`, cluesheetID, ipaUID)
	if err != nil {
		return model.UserParticipation{}, fmt.Errorf("get participation: %w", err)
	}

	participation, err := pgx.CollectOneRow[model.UserParticipation](rows, pgx.RowToStructByNameLax[model.UserParticipation])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.UserParticipation{Cluesheet_id: cluesheetID, Ipa_uid: ipaUID, Hidden: false}, nil
	} else if err != nil {
		return model.UserParticipation{}, fmt.Errorf("get participation: %w", err)
	}

	return participation, nil
}

func UpsertParticipation(ctx context.Context, conn *pgxpool.Pool, participation model.UserParticipation) error {
	_, err := conn.Exec(ctx, `insert into user_participation(cluesheet_id, ipa_uid, hidden) values ($1, $2, $3) on conflict (cluesheet_id, ipa_uid) do update set hidden = $3`, participation.Cluesheet_id, participation.Ipa_uid, participation.Hidden)
	if err != nil {
		return fmt.Errorf("upsert participation: %w", err)
	}

	return nil
}
