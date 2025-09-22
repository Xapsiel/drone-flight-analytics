package migration

import (
	"embed"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(pgURL string) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		slog.Error("error loading migration files", "err", err)
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, pgURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
