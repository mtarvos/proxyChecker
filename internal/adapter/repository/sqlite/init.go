package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(storagePath string, log *slog.Logger) (*Storage, error) {
	const fn = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s opening sqlite db: %w", fn, err)
	}

	storage := &Storage{db: db, log: log}

	err = storage.MigrationsUP()
	if err != nil {
		return nil, fmt.Errorf("%s MigrationsUP failed: %e", fn, err)
	}

	return storage, nil
}
