package services

import (
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type TeamService interface {
	Add(team dto.Team) error
	Get(teamName string) (dto.Team, error)
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{}
}

func (t *teamService) Add(team dto.Team) error {
	return nil
}

func (t *teamService) Get(teamName string) (dto.Team, error) {
	return dto.Team{}, nil
}
