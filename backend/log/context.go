package log

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"go.uber.org/zap"

	"csh/cluesheet/config"
)

type logContextKey string

const contextKey logContextKey = "log"

/**
* Instantiate a new logger
* Be sure to defer logger.Sync() to ensure it's flushed before we exit
 */
func GetLogger(ctx context.Context) *zap.Logger {
	// TODO we should probably handle the error
	var logger *zap.Logger

	switch config.FromContext(ctx).GetString("environment") {
	case "local":
		logger, _ = zap.NewDevelopment()
	case "dev":
		logger, _ = zap.NewDevelopment()
	case "prod":
		logger, _ = zap.NewProduction()
	}

	return logger
}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	l := ctx.Value(contextKey)
	var logger *zap.Logger
	if l == nil {
		logger = zap.L()
	} else {
		logger = l.(*zap.Logger)
	}

	config := config.FromContext(ctx)
	if span, ok := tracer.SpanFromContext(ctx); ok {
		spanContext := span.Context()
		logger = logger.With(
			zap.String("dd.trace_id", spanContext.TraceID()),
			zap.Uint64("dd.span_id", spanContext.SpanID()),
			zap.String("dd.version", config.GetString("version")),
			zap.String("dd.env", config.GetString("environment")),
		)
	}
	return logger
}
