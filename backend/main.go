package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		rw.Write([]byte(r.RemoteAddr))
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
