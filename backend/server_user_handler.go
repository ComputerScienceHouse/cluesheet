package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func handleGetUserCluesheet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cluesheet_id, err := uuid.Parse(vars["cluesheet_id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse cluesheet uuid '%s': '%s'", vars["cluesheet_id"], err))
		return
	}

	ipa_uid := vars["ipa_uid"]

	rows, err := conn.Query(r.Context(), `select id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups from cluesheet where id = $1`, cluesheet_id)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed getting cluesheet '%s': '%s'", cluesheet_id, err))
		return
	}

	cluesheet, err := pgx.CollectOneRow[Cluesheet](rows, pgx.RowToStructByNameLax[Cluesheet])
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed getting cluesheet '%s': '%s'", cluesheet_id, err))
		return
	}

	clues, err := GetCluesForUser(r.Context(), conn, cluesheet.Id, ipa_uid)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed resolving clues for user '%s' on cluesheet '%s': '%s'", ipa_uid, cluesheet_id, err))
		return
	}
	cluesheet.Clues = &clues

	writeJSON(rw, http.StatusOK, cluesheet)
}
