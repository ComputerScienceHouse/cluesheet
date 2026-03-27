package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"csh/cluesheet/config"
	"csh/cluesheet/log"

	muxtrace "github.com/DataDog/dd-trace-go/contrib/gorilla/mux/v2"
	pgxtrace "github.com/DataDog/dd-trace-go/contrib/jackc/pgx.v5/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"

	"go.uber.org/zap"
)

const (
	connStr = "postgres://postgres:test@localhost:5432/postgres?sslmode=disable"
)

func main() {
	ctx := context.Background()

	ctx = config.ContextWithConfig(ctx, config.GetConfig(ctx))
	ctx = log.ContextWithLogger(ctx, log.GetLogger(ctx))

	// Tag all logs with the version string if we have one
	if config.FromContext(ctx).GetString("version") != "" {
		ctx = log.ContextWithLogger(ctx, log.GetLogger(ctx).With(
			zap.String("version", config.FromContext(ctx).GetString("version")),
		))
	}

	if config.FromContext(ctx).GetBool("tracing.enabled") {
		tracer.Start(
			tracer.WithEnv(config.FromContext(ctx).GetString("env")),
			tracer.WithService("cluesheet"),
			tracer.WithServiceVersion(config.FromContext(ctx).GetString("version")),
		)
		defer tracer.Stop()
		log.FromContext(ctx).Debug("started tracing")
	}

	conn, err := pgxtrace.NewPool(ctx, connStr)
	if err != nil {
		log.FromContext(ctx).Fatal("failed setting up db pool", zap.Error(err))
	}
	defer conn.Close()

	router := muxtrace.NewRouter()

	registerRoutes(ctx, router, conn)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	// channel to prevent the main thread from exiting too early
	shutdown := make(chan struct{})

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
		log.FromContext(ctx).Info("shutting down gracefully")
		srv.Shutdown(ctx)
		shutdown <- struct{}{}
	}()

	log.FromContext(ctx).Info("starting http server", zap.String("addr", srv.Addr))
	srv.ListenAndServe()
	<-shutdown
	log.FromContext(ctx).Info("exiting...")
}
