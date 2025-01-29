package sqlite

import (
	"database/sql"
	"log/slog"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}
