package app

import (
	"fmt"
	"log/slog"
	someboundedcontext "project_template/internal/someboundedcontext/config"
	"project_template/pkg/database"
	"project_template/pkg/webserver"
	"strings"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	WebServer          webserver.Config
	SomeBoundedContext someboundedcontext.Config
	Database           database.Config
	fx.Out
}

func NewServeConfig(yamlConfigFile string) func() (Config, error) {
	return func() (Config, error) {
		v := viper.New()

		// Config file settings
		if yamlConfigFile != "" {
			v.SetConfigFile(yamlConfigFile)

			if err := v.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					return Config{}, fmt.Errorf("reading config file: %w", err)
				}
			}
		} else {
			slog.Warn("No config file provided")
		}

		// Environment variables
		v.SetEnvPrefix("APP")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		v.AutomaticEnv()

		// Unmarshal to struct
		var cfg Config
		if err := defaults.Set(&cfg); err != nil {
			return Config{}, fmt.Errorf("setting defaults: %w", err)
		}
		if err := v.Unmarshal(&cfg); err != nil {
			return Config{}, fmt.Errorf("parsing config: %w", err)
		}

		return cfg, nil
	}
}
