package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/anish-chanda/ferna/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed postgresql/*.sql
var PostgresMigrations embed.FS

// RunMigrations executes database migrations using the embedded PostgreSQL files
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, logger *logger.ServiceLogger) error {
	logger.Info("Starting PostgreSQL database migrations")

	// Create source driver from embedded filesystem
	source, err := iofs.New(PostgresMigrations, "postgresql")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	defer source.Close()

	// Convert pgxpool to sql.DB for migrate
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	// Create database driver for postgres
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current migration version
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is in a dirty state at version %d", currentVersion)
	}

	if err == migrate.ErrNilVersion {
		logger.Debug("No existing migrations found, starting fresh")
	} else {
		logger.Debugf("Current database version: %d", currentVersion)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		logger.Debug("No new migrations to apply")
	} else {
		// Get new version after migration
		newVersion, _, err := m.Version()
		if err != nil {
			return fmt.Errorf("failed to get new migration version: %w", err)
		}
		logger.Infof("Migrations completed successfully. Database version: %d", newVersion)
	}

	return nil
}
