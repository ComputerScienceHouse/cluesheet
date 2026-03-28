package dbcluesheet

import (
	"context"
	"fmt"

	"csh/cluesheet/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListCluesheets(ctx context.Context, conn *pgxpool.Pool) ([]model.Cluesheet, error) {
	rows, err := conn.Query(ctx, `select id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups from cluesheet`)
	if err != nil {
		return nil, fmt.Errorf("list cluesheets: %w", err)
	}

	cluesheets, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Cluesheet])
	if err != nil {
		return nil, fmt.Errorf("list cluesheets: %w", err)
	}

	return cluesheets, nil
}

func GetCluesheet(ctx context.Context, conn *pgxpool.Pool, id uuid.UUID) (model.Cluesheet, error) {
	rows, err := conn.Query(ctx, `select id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups from cluesheet where id = $1`, id)
	if err != nil {
		return model.Cluesheet{}, fmt.Errorf("get cluesheet: %w", err)
	}

	cluesheet, err := pgx.CollectOneRow[model.Cluesheet](rows, pgx.RowToStructByNameLax[model.Cluesheet])
	if err != nil {
		return model.Cluesheet{}, fmt.Errorf("get cluesheet: %w", err)
	}

	return cluesheet, nil
}

func CreateCluesheet(ctx context.Context, conn *pgxpool.Pool, sheet model.Cluesheet) error {
	_, err := conn.Exec(ctx, `insert into cluesheet (id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, sheet.Id, sheet.Name, sheet.Origin_id, sheet.Created_by, sheet.Created_at, sheet.Edited_by, sheet.Edited_at, sheet.Visibility, sheet.Owners, sheet.Groups)
	if err != nil {
		return fmt.Errorf("create cluesheet: %w", err)
	}

	return nil
}
