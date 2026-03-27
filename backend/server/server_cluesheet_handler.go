package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"csh/cluesheet/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func handleListCluesheets(rw http.ResponseWriter, r *http.Request) {
	rows, err := conn.Query(r.Context(), `select id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups from cluesheet`)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed getting cluesheets: '%s'", err))
		return
	}

	cluesheets, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Cluesheet])
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed getting cluesheet: '%s'", err))
		return
	}

	for _, cluesheet := range cluesheets {
		clues, err := model.GetClues(r.Context(), conn, cluesheet.Id)
		if err != nil {
			writeError(rw, 500, fmt.Sprintf("failed resolving clues: '%s'", err))
			return
		}
		cluesheet.Clues = &clues
	}

	writeJSON(rw, http.StatusOK, cluesheets)
}

func handleGetCluesheet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse uuid '%s': '%s'", vars["id"], err))
		return
	}

	rows, err := conn.Query(r.Context(), `select id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups from cluesheet where id = $1`, id)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed getting cluesheet '%s': '%s'", vars["id"], err))
		return
	}

	cluesheet, err := pgx.CollectOneRow[model.Cluesheet](rows, pgx.RowToStructByNameLax[model.Cluesheet])
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed getting cluesheet '%s': '%s'", vars["id"], err))
		return
	}

	clues, err := model.GetClues(r.Context(), conn, cluesheet.Id)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed resolving clues '%s': '%s'", vars["id"], err))
		return
	}
	cluesheet.Clues = &clues

	writeJSON(rw, http.StatusOK, cluesheet)
}

/*
POST /api/v1/cluesheet: create a cluesheet
params:
  - ID (generated)
  - Name
  - Origin (optional)
  - creator
  - created time (generated)
  - editor (from creator)
  - edited time (generated)
  - visibility (optional, default hidden)
  - Owners (optional, default creator)
  - Groups (optional)
*/
func handlePostCluesheet(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to read request: %s", err.Error()))
		return
	}

	var params model.PostCluesheetParams

	err = json.Unmarshal(body, &params)
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse request: %s", err.Error()))
		return
	}
	if params.Name == "" {
		writeError(rw, 400, "Name must not be empty")
		return
	}
	if params.Creator == "" {
		writeError(rw, 400, "Creator must not be empty")
		return
	}

	if params.Visibility == nil || *params.Visibility == "" {
		// TODO this isn't canonical
		local := "hidden"
		params.Visibility = &local
	}

	if !slices.Contains(params.Owners, params.Creator) {
		params.Owners = append(params.Owners, params.Creator)
	}

	if params.Groups == nil {
		params.Groups = []string{}
	}

	newSheet := model.Cluesheet{
		Id:         uuid.New(),
		Name:       params.Name,
		Origin_id:  params.Origin,
		Created_by: params.Creator,
		Created_at: time.Now(),
		Edited_by:  params.Creator,
		Edited_at:  time.Now(),
		Visibility: *params.Visibility,
		Owners:     params.Owners,
		Groups:     params.Groups,
	}
	_, err = conn.Exec(r.Context(), `insert into cluesheet (id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, newSheet.Id, newSheet.Name, newSheet.Origin_id, newSheet.Created_by, newSheet.Created_at, newSheet.Edited_by, newSheet.Edited_at, newSheet.Visibility, newSheet.Owners, newSheet.Groups)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to persist cluesheet: %s", err))
		return
	}

	writeJSON(rw, http.StatusOK, newSheet)
}
