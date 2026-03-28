package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	dbprogress "csh/cluesheet/db/progress"
	"csh/cluesheet/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	progress, err := dbprogress.GetProgress(r.Context(), conn, clue_id, ipa_uid)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to query for user progress on clue '%s' for user '%s': %s", clue_id, ipa_uid, err))
		return
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

	progress := model.UserProgress{
		Ipa_uid:     ipa_uid,
		Clue_id:     clue_id,
		Completions: completions,
	}
	if err := dbprogress.UpsertProgress(r.Context(), conn, progress); err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to update user progress on clue '%s' for user '%s' to value '%d': %s", clue_id, ipa_uid, completions, err))
		return
	}
	writeJSON(rw, http.StatusOK, progress)
}
