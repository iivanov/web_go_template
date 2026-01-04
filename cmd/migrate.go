package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"project_template/pkg/database"
	"project_template/pkg/migrations"

	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Run database migrations using Goose. Supports up, down, and status subcommands.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator, err := createMigrator()
		if err != nil {
			return err
		}
		return migrator.Up()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator, err := createMigrator()
		if err != nil {
			return err
		}
		return migrator.Down()
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator, err := createMigrator()
		if err != nil {
			return err
		}
		return migrator.Status()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}

func createMigrator() (*migrations.Migrator, error) {
	if err := godotenv.Load(); err == nil {
		slog.Info("Loaded .env file")
	}

	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("reading config file: %w", err)
			}
		}
	}

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	var dbConfig database.Config
	if err := defaults.Set(&dbConfig); err != nil {
		return nil, fmt.Errorf("setting defaults: %w", err)
	}

	if err := v.UnmarshalKey("database", &dbConfig); err != nil {
		return nil, fmt.Errorf("parsing database config: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.NewConnection(logger, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting sql.DB: %w", err)
	}

	return migrations.NewMigrator(sqlDB, logger, migrations.Config{}), nil
}
