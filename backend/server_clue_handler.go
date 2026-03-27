package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

	newClue := Clue{
		Id:          uuid.New(),
		Description: params.Description,
		Tags:        params.Tags,
		Origin_id:   id,
		Created_by:  params.Creator,
		Created_at:  time.Now(),
		Edited_by:   params.Creator,
		Edited_at:   time.Now(),
		Children:    []*Clue{},
	}

	tx, err := conn.Begin(r.Context())
	if err != nil {
		writeError(rw, 500, "failed to store clue")
		fmt.Println(err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	_, err = tx.Exec(r.Context(), `insert into clue(id, description, tags, origin_id, created_by, created_at, edited_by, edited_at) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		newClue.Id,
		newClue.Description,
		newClue.Tags,
		newClue.Origin_id,
		newClue.Created_by,
		newClue.Created_at,
		newClue.Edited_by,
		newClue.Edited_at,
	)
	if err != nil {
		writeError(rw, 500, "failed to store clue")
		fmt.Println(err.Error())
		return
	}

	var cr *ClueRelation = nil

	// fmt.Printf("%v\n%s\n", params, body) // printing was here presumably for debugging

	if params.ParentClueId != nil {
		cr = &ClueRelation{
			Id:        uuid.New(),
			Parent_id: *params.ParentClueId,
			Child_id:  newClue.Id,
		}

		// TODO validate parent exists on same sheet

		_, err = tx.Exec(r.Context(), `insert into clue_relation(id, parent_id, child_id) values ($1, $2, $3)`, cr.Id, cr.Parent_id, cr.Child_id)
		if err != nil {
			writeError(rw, 500, "failed to store clue parent")
			fmt.Println(err.Error())
			return
		}
	}

	err = tx.Commit(r.Context())
	if err != nil {
		writeError(rw, 500, "failed to store clue with parent")
		fmt.Println(err.Error())
		return
	}

	writeJSON(rw, http.StatusOK, struct {
		Clue         Clue
		ClueRelation *ClueRelation `json:",omitempty,omitzero"`
	}{Clue: newClue, ClueRelation: cr})
}
