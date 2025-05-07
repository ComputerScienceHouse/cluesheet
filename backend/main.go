package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"time"

	"csh/cluesheet/config"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"

	muxtrace "github.com/DataDog/dd-trace-go/contrib/gorilla/mux/v2"
	"github.com/gorilla/mux"

	pgxtrace "github.com/DataDog/dd-trace-go/contrib/jackc/pgx.v5/v2"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

const (
	connStr = "postgres://postgres:test@localhost:5432/postgres?sslmode=disable"
)

func main() {
	ctx := context.Background()

	ctx = config.ContextWithConfig(ctx, config.GetConfig(ctx))

	if config.FromContext(ctx).GetBool("tracing.enabled") {
		tracer.Start(
			tracer.WithEnv(config.FromContext(ctx).GetString("env")),
			tracer.WithService("cluesheet"),
			// tracer.WithServiceVersion(), // TODO add once we have a git commit
		)
		defer tracer.Stop()
	}

	conn, err := pgxtrace.NewPool(ctx, connStr)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	router := muxtrace.NewRouter()

	// Pass the config in context to all requests
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = config.ContextWithConfig(ctx, config.GetConfig(ctx))
			r = r.WithContext(ctx)
			h.ServeHTTP(rw, r)
		})
	})

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(r.RemoteAddr))
	})

	v1 := router.PathPrefix("/api/v1/").Subrouter()
	v1.Path("/cluesheet").Methods("GET").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rows, err := conn.Query(r.Context(), `select id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups from cluesheet`)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed getting cluesheets: '%s'", err), 500)
			return
		}

		cluesheets, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Cluesheet])
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed getting cluesheet: '%s'", err), 500)
			return
		}

		for _, cluesheet := range cluesheets {
			clues, err := GetClues(ctx, conn, cluesheet.Id)
			if err != nil {
				http.Error(rw, fmt.Sprintf("failed resolving clues: '%s'", err), 500)
				return
			}
			cluesheet.Clues = &clues
		}

		data, err := json.Marshal(cluesheets)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed marshalling cluesheets: '%s'", err), 500)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(data)
	})
	v1.Path("/cluesheet/{id}").Methods("GET").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := uuid.Parse(vars["id"])
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to parse uuid '%s': '%s'", vars["id"], err), 400)
			return
		}

		rows, err := conn.Query(r.Context(), `select id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups from cluesheet where id = $1`, id)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed getting cluesheet '%s': '%s'", vars["id"], err), 500)
			return
		}

		cluesheet, err := pgx.CollectOneRow[Cluesheet](rows, pgx.RowToStructByNameLax[Cluesheet])
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed getting cluesheet '%s': '%s'", vars["id"], err), 500)
			return
		}

		clues, err := GetClues(ctx, conn, cluesheet.Id)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed resolving clues '%s': '%s'", vars["id"], err), 500)
			return
		}
		cluesheet.Clues = &clues

		data, err := json.Marshal(cluesheet)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed marshalling cluesheet '%s': '%s'", vars["id"], err), 500)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(data)
	})

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
	v1.Path("/cluesheet").Methods("POST").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to read request: %s", err.Error()), 400)
			return
		}

		var params struct {
			Name       string
			Origin     *uuid.UUID
			Creator    string
			Visibility *string
			Owners     []string
			Groups     []string
		}

		err = json.Unmarshal(body, &params)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to parse request: %s", err.Error()), 400)
			return
		}
		if params.Name == "" {
			http.Error(rw, "Name must not be empty", 400)
			return
		}
		if params.Creator == "" {
			http.Error(rw, "Creator must not be empty", 400)
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

		newSheet := Cluesheet{
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
		_, err = conn.Exec(ctx, `insert into cluesheet (id, name, origin_id, created_by, created_at, edited_by, edited_at, visibility, owners, groups) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, newSheet.Id, newSheet.Name, newSheet.Origin_id, newSheet.Created_by, newSheet.Created_at, newSheet.Edited_by, newSheet.Edited_at, newSheet.Visibility, newSheet.Owners, newSheet.Groups)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to persist cluesheet: %s", err), 500)
			return
		}

		data, err := json.Marshal(newSheet)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed marshalling cluesheet: '%s'", err), 500)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(data)
	})

	// POST /api/v1/cluesheet/{id}/clue
	// params:
	// - description
	// - rules are TBD
	// - creator
	// - created and edited are derived/generated
	// - tags (optional)
	// - parent_clue_id (optional)
	v1.Path("/cluesheet/{id}/clue").Methods("POST").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := uuid.Parse(vars["id"])
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to parse uuid '%s': '%s'", vars["id"], err), 400)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to read request: %s", err.Error()), 400)
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
			http.Error(rw, fmt.Sprintf("failed to parse request: %s", err.Error()), 400)
			return
		}

		if params.Description == "" {
			http.Error(rw, "Description must not be empty", 400)
			return
		}

		if params.Creator == "" {
			http.Error(rw, "Creator must not be empty", 400)
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
			http.Error(rw, "failed to store clue", 500)
			fmt.Printf(err.Error())
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
			http.Error(rw, "failed to store clue", 500)
			fmt.Printf(err.Error())
			return
		}

		var cr *ClueRelation = nil

		fmt.Printf("%v\n%s\n", params, body)

		if params.ParentClueId != nil {
			cr = &ClueRelation{
				Id:        uuid.New(),
				Parent_id: *params.ParentClueId,
				Child_id:  newClue.Id,
			}

			// TODO validate parent exists on same sheet

			_, err = tx.Exec(r.Context(), `insert into clue_relation(id, parent_id, child_id) values ($1, $2, $3)`, cr.Id, cr.Parent_id, cr.Child_id)
			if err != nil {
				http.Error(rw, "failed to store clue parent", 500)
				fmt.Printf(err.Error())
				return
			}
		}

		err = tx.Commit(r.Context())
		if err != nil {
			http.Error(rw, "failed to store clue with parent", 500)
			fmt.Printf(err.Error())
			return
		}

		data, err := json.Marshal(struct {
			Clue         Clue
			ClueRelation *ClueRelation `json:,omitempty,omitzero`
		}{Clue: newClue, ClueRelation: cr})
		if err != nil {
			http.Error(rw, "Failed to marshal data", 500)
			fmt.Printf(err.Error())
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(data)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancel := context.WithTimeout(ctx, time.Second*15)
		defer cancel()
		go func() {
			select {
			case <-c:
				// Immediately stop on a repeated sigterm
				cancel()
			case <-ctx.Done():
			}
		}()
		srv.Shutdown(ctx)
	}()

	srv.ListenAndServe()
}
