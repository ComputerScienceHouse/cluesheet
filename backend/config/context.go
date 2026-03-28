package config

import (
	"context"

	"github.com/spf13/viper"
)

type configContextKey string

const contextKey configContextKey = "config"

func ContextWithConfig(ctx context.Context, config *viper.Viper) context.Context {
	return context.WithValue(ctx, contextKey, config)
}

func FromContext(ctx context.Context) *viper.Viper {
	v := ctx.Value(contextKey)
	if v == nil {
		panic("config not available in context, somewhere it got lost")
	}
	return v.(*viper.Viper)
}
