package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	dbclue "csh/cluesheet/db/clue"
	"csh/cluesheet/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// POST /api/v1/cluesheet/{id}/clue
// params:
// - description
// - rules are TBD
// - creator
// - created and edited are derived/generated
// - tags (optional)
// - parent_clue_id (optional)
func handlePostClue(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse uuid '%s': '%s'", vars["id"], err))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to read request: %s", err.Error()))
		return
	}

	var params struct {
		Description  string
		Creator      string
		Tags         []string
		ParentClueId *uuid.UUID `json:"parent_clue_id"`
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse request: %s", err.Error()))
		return
	}

	if params.Description == "" {
		writeError(rw, 400, "Description must not be empty")
		return
	}

	if params.Creator == "" {
		writeError(rw, 400, "Creator must not be empty")
		return
	}

	if params.Tags == nil {
		params.Tags = []string{}
	}

	newClue := model.Clue{
		Id:          uuid.New(),
		Description: params.Description,
		Tags:        params.Tags,
		Origin_id:   id,
		Created_by:  params.Creator,
		Created_at:  time.Now(),
		Edited_by:   params.Creator,
		Edited_at:   time.Now(),
		Children:    []*model.Clue{},
	}

	result, err := dbclue.CreateClue(r.Context(), conn, newClue, params.ParentClueId)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to store clue: %s", err))
		return
	}
	writeJSON(rw, http.StatusOK, result)
}
