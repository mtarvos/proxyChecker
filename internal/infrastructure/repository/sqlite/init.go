package sqlite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Storage struct {
	db   *sqlx.DB
	log  *slog.Logger
	used map[string]map[string]bool
}

func New(storagePath string, log *slog.Logger) (*Storage, error) {
	const fn = "sqlite.New"

	db, err := sqlx.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s opening sqlite db: %w", fn, err)
	}

	uniq := make(map[string]map[string]bool)

	storage := &Storage{db: db, log: log, used: uniq}

	err = storage.MigrationsUP()
	if err != nil {
		return nil, fmt.Errorf("%s MigrationsUP failed: %e", fn, err)
	}

	return storage, nil
}
