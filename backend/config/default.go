package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

/*
* GetConfig determines what config file to look at based on the environment, and constructs a config to use
 */
func GetConfig(_ context.Context) *viper.Viper {
	config := viper.New()
	config.SetEnvPrefix("cluesheet")

	config.BindEnv("config_path")
	config.BindEnv("environment")
	config.SetDefault("environment", "local")

	if config.IsSet("config_path") {
		config.SetConfigFile(config.GetString("config_path"))
	} else {
		config.SetConfigType("yaml")
		config.AddConfigPath(".")
		config.AddConfigPath("/opt/cluesheet")

		switch config.GetString("environment") {
		case "local":
			config.SetConfigName("config.yaml")
		case "develop":
			config.SetConfigName("config.dev")
		case "prod":
			config.SetConfigName("config.prod")
		default:
			panic(fmt.Sprintf("unknown environment '%s', can't look for config", config.GetString("environment")))
		}
	}
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf(err.Error()))
	}

	return config
}
