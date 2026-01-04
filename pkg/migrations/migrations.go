package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

type Migrator struct {
	db     *sql.DB
	logger *slog.Logger
	config Config
}

func NewMigrator(db *sql.DB, logger *slog.Logger, config Config) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
		config: config,
	}
}

func (m *Migrator) Up() error {
	goose.SetLogger(goose.NopLogger())

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(m.db, "sql"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("database migrations completed successfully")
	return nil
}

func (m *Migrator) Down() error {
	goose.SetLogger(goose.NopLogger())

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Down(m.db, "sql"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	m.logger.Info("database migration rolled back successfully")
	return nil
}

func (m *Migrator) Status() error {
	goose.SetLogger(goose.NopLogger())

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Status(m.db, "sql"); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}
