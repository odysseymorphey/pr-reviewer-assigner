package services

import (
	"context"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type TeamService interface {
	Add(ctx context.Context, team dto.Team) error
	Get(teamName string) ([]dto.TeamMember, error)
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{
		repo: repo,
	}
}

func (s *teamService) Add(ctx context.Context, team dto.Team) error {
	return s.repo.Add(ctx, team)
}

func (s *teamService) Get(teamName string) ([]dto.TeamMember, error) {
	return s.repo.Get(teamName)
}
