package server

import (
	"context"
	"encoding/json"
	"net/http"

	"csh/cluesheet/config"
	"csh/cluesheet/log"

	muxtrace "github.com/DataDog/dd-trace-go/contrib/gorilla/mux/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

/**
* RegisterRoutes applies all the routes to the given mux router
* It's intended to break things out of main for testing
*
* rootContext is the context that the router's logger and config will be read
*     from, it is expected to be descended from the context in main
*
* router is the mux router to register routes on
*
* srvConn is the database connection
 */
func RegisterRoutes(rootContext context.Context, router *muxtrace.Router, srvConn *pgxpool.Pool) {
	conn = srvConn

	// Pass the logger and config in context to all requests
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = config.ContextWithConfig(ctx, config.FromContext(rootContext))
			ctx = log.ContextWithLogger(ctx, log.FromContext(rootContext))
			r = r.WithContext(ctx)
			h.ServeHTTP(rw, r)
		})
	})

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(r.RemoteAddr))
	})

	v1 := router.PathPrefix("/api/v1/").Subrouter()
	v1.Path("/cluesheet").Methods("GET").HandlerFunc(handleListCluesheets)
	v1.Path("/cluesheet/{id}").Methods("GET").HandlerFunc(handleGetCluesheet)
	v1.Path("/cluesheet").Methods("POST").HandlerFunc(handlePostCluesheet)
	v1.Path("/cluesheet/{id}/clue").Methods("POST").HandlerFunc(handlePostClue)
	v1.Path("/cluesheet/{cluesheet_id}/clue/{clue_id}/progress/{ipa_uid}").Methods("GET").HandlerFunc(handleGetProgress)
	v1.Path("/cluesheet/{cluesheet_id}/clue/{clue_id}/progress/{ipa_uid}").Methods("POST").HandlerFunc(handlePostProgress)
	v1.Path("/cluesheet/{cluesheet_id}/participation/{ipa_uid}").Methods("GET").HandlerFunc(handleGetParticipation)
	v1.Path("/cluesheet/{cluesheet_id}/participation/{ipa_uid}").Methods("POST").HandlerFunc(handlePostParticipation)
	v1.Path("/cluesheet/{cluesheet_id}/user/{ipa_uid}").Methods("GET").HandlerFunc(handleGetUserCluesheet)
}

func writeJSON(rw http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(data)
}

func writeError(rw http.ResponseWriter, status int, msg string) {
	writeJSON(rw, status, struct {
		Error string `json:"error"`
	}{Error: msg})
}
