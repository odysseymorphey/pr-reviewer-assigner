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

func (r *teamRepo) Get(teamName string) (dto.Team, error) {
	const query = ``

	return dto.Team{}, nil
}

func (r *teamRepo) Add(team dto.Team) error {
	return nil
}
