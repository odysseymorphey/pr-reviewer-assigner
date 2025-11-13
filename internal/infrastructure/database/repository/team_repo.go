package repository

import (
	"database/sql"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type teamRepo struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) repository.TeamRepository {
	return &teamRepo{}
}

func (t *teamRepo) Get() (dto.Team, error) {
	return dto.Team{}, nil
}

func (t *teamRepo) Add(team dto.Team) error {
	return nil
}
