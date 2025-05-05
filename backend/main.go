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

	"github.com/gorilla/mux"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
)

const (
	connStr = "postgres://postgres:test@localhost:5432/postgres?sslmode=disable"
)

func main() {
	ctx := context.Background()

	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	_, err = conn.Exec(ctx, `insert into cluesheet (id, created_at) values
	($1, $2)`, uuid.New(), time.Now())
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		rw.Write([]byte(r.RemoteAddr))
	})

	rows, err := conn.Query(ctx, `select id, created_at from cluesheet`)
	if err != nil {
		panic(err.Error())
	}
	/*
		// For simpler error handling, consider using the higher-level pgx v5
		// CollectRows() and ForEachRow() helpers instead.

	*/
	if rows.Err() != nil {
		panic(rows.Err().Error())
	}

	for rows.Next() {
		var id *uuid.UUID
		var creation *time.Time
		err := rows.Scan(&id, &creation)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%v, %v\n", id, creation)
	}

	v1 := router.PathPrefix("/api/v1/").Subrouter()
	v1.Path("/cluesheet/{id}").Methods("GET").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
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
