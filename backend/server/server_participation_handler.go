package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"csh/cluesheet/log"
	"csh/cluesheet/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func handleGetParticipation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cluesheet_id, err := uuid.Parse(vars["cluesheet_id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse uuid '%s': '%s'", vars["cluesheet_id"], err))
		return
	}

	ipa_uid := vars["ipa_uid"]
	// todo validation

	rows, err := conn.Query(r.Context(), `select * from user_participation where cluesheet_id = $1 and ipa_uid = $2`, cluesheet_id, ipa_uid)
	if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to query for user participation on '%s' for user '%s': %s", cluesheet_id, ipa_uid, err))
		return
	}

	participation, err := pgx.CollectOneRow[model.UserParticipation](rows, pgx.RowToStructByNameLax[model.UserParticipation])
	if err == pgx.ErrNoRows {
		// if no stored result, there's no hiding
		participation = model.UserParticipation{
			Cluesheet_id: cluesheet_id,
			Ipa_uid:      ipa_uid,
			Hidden:       false,
		}
	} else if err != nil {
		writeError(rw, 500, fmt.Sprintf("failed to query for user participation on '%s' for user '%s': %s", cluesheet_id, ipa_uid, err))
		return
	}

	writeJSON(rw, http.StatusOK, participation)
}

func handlePostParticipation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cluesheet_id, err := uuid.Parse(vars["cluesheet_id"])
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse uuid '%s': '%s'", vars["cluesheet_id"], err))
		return
	}

	ipa_uid := vars["ipa_uid"]
	// todo validation

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to read request: %s", err.Error()))
		return
	}

	var params struct {
		Hidden bool
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		writeError(rw, 400, fmt.Sprintf("failed to parse request: %s", err.Error()))
		return
	}

	participation := model.UserParticipation{Cluesheet_id: cluesheet_id, Ipa_uid: ipa_uid, Hidden: params.Hidden}

	log.FromContext(r.Context()).Info("participation", zap.Any("participation", participation))
	_, err = conn.Exec(r.Context(), `insert into user_participation(cluesheet_id, ipa_uid, hidden) values ($1, $2, $3) on conflict (cluesheet_id, ipa_uid) do update set hidden = $3`, participation.Cluesheet_id, participation.Ipa_uid, participation.Hidden)
	if err != nil {
		writeError(rw, 500, "failed to store clue parent")
		fmt.Println(err.Error())
		return
	}

	writeJSON(rw, http.StatusOK, participation)
}
