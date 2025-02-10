package sqlite

import (
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
)

var (
	//go:embed "migrations/*.sql"
	migrationsFS embed.FS
)

func (s *Storage) MigrationsUP() error {
	goose.SetBaseFS(migrationsFS)
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return fmt.Errorf("can not set up sqlite3 dialect: %w", err)
	}
	migrationsDir := "migrations"

	if err = goose.Up(s.db.DB, migrationsDir); err != nil {
		return fmt.Errorf("can not up migrations: %w")
	}

	return nil
}
