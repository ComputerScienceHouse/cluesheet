package server

import (
	"fmt"
	"net/http"

	dbclue "csh/cluesheet/db/clue"
	dbcluesheet "csh/cluesheet/db/cluesheet"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func handleGetUserCluesheet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cluesheet_id, err := uuid.Parse(vars["cluesheet_id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse cluesheet uuid '%s': '%s'", vars["cluesheet_id"], err))
		return
	}

	ipa_uid := vars["ipa_uid"]

	cluesheet, err := dbcluesheet.GetCluesheet(r.Context(), conn, cluesheet_id)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed getting cluesheet '%s': '%s'", cluesheet_id, err))
		return
	}
	clues, err := dbclue.GetCluesForUser(r.Context(), conn, cluesheet.Id, ipa_uid)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed resolving clues for user '%s' on cluesheet '%s': '%s'", ipa_uid, cluesheet_id, err))
		return
	}
	cluesheet.Clues = &clues
	writeJSON(rw, http.StatusOK, cluesheet)
}
