package sqlite

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func parseTime(value string) (time.Time, error) {
	return time.ParseInLocation(
		time.DateTime,
		value,
		time.UTC,
	)
}
