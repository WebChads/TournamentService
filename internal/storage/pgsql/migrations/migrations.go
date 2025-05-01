package migrations

import (
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(db *sql.DB, logger *slog.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("could not create migration driver", slogerr.Error(err))
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations", "postgres", driver,
	)
	if err != nil {
		logger.Error("could not create migration instanse", slogerr.Error(err))
		return err
	}

	// Apply available migrations
	if err := migration.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Warn("could not apply migrations", slogerr.Error(err))
			return nil
		}

		if strings.Contains(err.Error(), "no such table") {
			logger.Warn(
				"first-time migration, creating schema_migrations table",
				slogerr.Warn(err),
			)

			if err := migration.Force(1); err != nil {
				return err
			}

			return migration.Up()
		}

		logger.Error("could not apply migrations", slogerr.Error(err))
		return err
	}

	return nil
}
