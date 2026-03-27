package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	muxtrace "github.com/DataDog/dd-trace-go/contrib/gorilla/mux/v2"
	pgxtrace "github.com/DataDog/dd-trace-go/contrib/jackc/pgx.v5/v2"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	"csh/cluesheet/config"
)

func TestCluesheetSnapshot(t *testing.T) {
	ctx := t.Context()

	ctx = config.ContextWithConfig(ctx, config.GetConfig(ctx))

	conn, err := pgxtrace.NewPool(ctx, connStr)
	t.Cleanup(conn.Close)
	assert.Nil(t, err, "failed to connect to postgres")

	router := muxtrace.NewRouter()

	registerRoutes(ctx, router, conn)

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	var cluesheetId uuid.UUID

	params := PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	}

	t.Run("PostCluesheet", func(t *testing.T) {
		buf, err := json.Marshal(params)
		assert.Nil(t, err, "failed to marshal params")

		resp, err := http.Post(server.URL+"/api/v1/cluesheet", "application/json", bytes.NewReader(buf))
		assert.Nil(t, err, "failed to post cluesheet")

		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err, "failed to read body")

		var cluesheet Cluesheet
		err = json.Unmarshal(body, &cluesheet)
		assert.Nil(t, err, "failed to unmarshal body")

		assert.Equal(t, params.Name, cluesheet.Name)
		assert.Nil(t, cluesheet.Origin_id)
		assert.Equal(t, params.Creator, cluesheet.Created_by)
		assert.Equal(t, params.Owners, cluesheet.Owners)
		assert.Equal(t, params.Groups, cluesheet.Groups)
		assert.Equal(t, "hidden", cluesheet.Visibility)
		assert.Nil(t, cluesheet.Clues)

		cluesheetId = cluesheet.Id

		t.Log("created a cluesheet", cluesheet)
	})

	t.Run("GetCluesheet", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/cluesheet/" + cluesheetId.String())
		assert.Nil(t, err, "failed to get cluesheet")

		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err, "failed to read body")

		var cluesheet Cluesheet
		err = json.Unmarshal(body, &cluesheet)
		assert.Nil(t, err, "failed to unmarshal body")

		assert.Equal(t, cluesheetId, cluesheet.Id)
		assert.Equal(t, params.Name, cluesheet.Name)
		assert.Nil(t, cluesheet.Origin_id)
		assert.Equal(t, params.Creator, cluesheet.Created_by)
		assert.Equal(t, params.Owners, cluesheet.Owners)
		assert.Equal(t, params.Groups, cluesheet.Groups)
		assert.Equal(t, "hidden", cluesheet.Visibility)
		assert.Nil(t, cluesheet.Clues)

		t.Log("got a cluesheet", cluesheet)
	})
}
