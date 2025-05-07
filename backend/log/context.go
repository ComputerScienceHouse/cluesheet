package log

import (
	"context"

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
	if l == nil {
		return zap.L()
	}
	return l.(*zap.Logger)
}
