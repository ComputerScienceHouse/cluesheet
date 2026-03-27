package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func handleGetProgress(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// I don't think we actually care about the sheet ID
	/*
		cluesheet_id, err := uuid.Parse(vars["cluesheet_id"])
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to parse cluesheet uuid '%s': '%s'", vars["cluesheet_id"], err), 400)
			return
		}
	*/

	clue_id, err := uuid.Parse(vars["clue_id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse clue uuid '%s': '%s'", vars["clue_id"], err))
		return
	}

	ipa_uid := vars["ipa_uid"]
	// todo validation

	rows, err := conn.Query(r.Context(), `select * from user_progress where clue_id = $1 and ipa_uid = $2`, clue_id, ipa_uid)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to query for user progress on clue '%s' for user '%s': %s", clue_id, ipa_uid, err))
		return
	}

	progress, err := pgx.CollectOneRow[UserProgress](rows, pgx.RowToStructByNameLax[UserProgress])
	if err == pgx.ErrNoRows {
		progress = UserProgress{
			Ipa_uid:     ipa_uid,
			Clue_id:     clue_id,
			Completions: 0,
		}
	} else if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to query for user progress on clue '%s' for user '%s': %s", clue_id, ipa_uid, err))
	}

	writeJSON(rw, http.StatusOK, progress)
}

// Body as a plain number
func handlePostProgress(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// I don't think we actually care about the sheet ID
	/*
		cluesheet_id, err := uuid.Parse(vars["cluesheet_id"])
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to parse cluesheet uuid '%s': '%s'", vars["cluesheet_id"], err), 400)
			return
		}
	*/
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to read request: %s", err.Error()))
		return
	}

	// do we care to be stricter in the parsing?
	completions, err := strconv.Atoi(string(body))
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse request: %s", err.Error()))
		return
	}

	clue_id, err := uuid.Parse(vars["clue_id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse clue uuid '%s': '%s'", vars["clue_id"], err))
		return
	}

	ipa_uid := vars["ipa_uid"]
	// todo validation

	_, err = conn.Exec(r.Context(), `insert into user_progress (ipa_uid, clue_id, completions) values ($1, $2, $3) on conflict (ipa_uid, clue_id) do update set completions = $3`, ipa_uid, clue_id, completions)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to update user progress on clue '%s' for user '%s' to value '%d': %s", clue_id, ipa_uid, completions, err))
		return
	}

	progress := UserProgress{
		Ipa_uid:     ipa_uid,
		Clue_id:     clue_id,
		Completions: completions,
	}

	writeJSON(rw, http.StatusOK, progress)
}
