package database

import (
	"database/sql"
	"pr-reviwer-assigner/internal/config"

	_ "github.com/lib/pq"
)

func New(cfg config.DBConfig) (*sql.DB, error) {
	pool, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
