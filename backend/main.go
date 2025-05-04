package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

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
