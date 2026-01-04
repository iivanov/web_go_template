package app

import (
	"fmt"
	"log/slog"
	someboundedcontext "project_template/internal/someboundedcontext/config"
	"project_template/pkg/database"
	"project_template/pkg/messagebus"
	"project_template/pkg/migrations"
	"project_template/pkg/telemetry"
	"project_template/pkg/webserver"
	"strings"

	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	WebServer          webserver.Config          `mapstructure:"webserver"`
	SomeBoundedContext someboundedcontext.Config `mapstructure:"someboundedcontext"`
	Database           database.Config           `mapstructure:"database"`
	Migrations         migrations.Config         `mapstructure:"migrations"`
	Telemetry          telemetry.Config          `mapstructure:"telemetry"`
	MessageBus         messagebus.Config         `mapstructure:"messagebus"`
}

func NewServeConfig(yamlConfigFile string) func() (Config, error) {
	return func() (Config, error) {
		// Load .env file if it exists
		if err := godotenv.Load(); err == nil {
			slog.Info("Loaded .env file")
		}

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

		// Bind telemetry config keys explicitly
		_ = v.BindEnv("telemetry.enabled", "APP_TELEMETRY_ENABLED")
		_ = v.BindEnv("telemetry.service_name", "APP_TELEMETRY_SERVICE_NAME")
		_ = v.BindEnv("telemetry.otlp_endpoint", "APP_TELEMETRY_OTLP_ENDPOINT")
		_ = v.BindEnv("telemetry.insecure", "APP_TELEMETRY_INSECURE")

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
