package migrations

import (
	"log/slog"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module("migrations",
	fx.Provide(func(db *gorm.DB, logger *slog.Logger, config Config) (*Migrator, error) {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		return NewMigrator(sqlDB, logger, config), nil
	}),
	fx.Invoke(RunMigrations),
)

func RunMigrations(m *Migrator) error {
	if !m.config.RunOnStartup {
		m.logger.Info("database migrations skipped (run_on_startup=false)")
		return nil
	}
	return m.Up()
}
