package repository

import (
	"database/sql"
	"pr-reviwer-assigner/internal/domain/repository"
)

type prRepo struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) repository.PRRepo {
	return &prRepo{
		db: db,
	}
}

func (r *prRepo) Create() {

}

func (r *prRepo) Merge() {

}

func (r *prRepo) Reassign() {

}
