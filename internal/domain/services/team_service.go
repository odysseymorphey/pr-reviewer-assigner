package services

import (
	"context"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type TeamService interface {
	Add(ctx context.Context, team dto.Team) error
	Get(teamName string) ([]dto.TeamMember, error)
	DeactivateMembers(ctx context.Context, req dto.TeamDeactivateRequest) (*dto.TeamDeactivateResponse, error)
}

type teamService struct {
	repo   repository.TeamRepository
	prRepo repository.PRRepository
}

func NewTeamService(repo repository.TeamRepository, prRepo repository.PRRepository) TeamService {
	return &teamService{
		repo:   repo,
		prRepo: prRepo,
	}
}

func (s *teamService) Add(ctx context.Context, team dto.Team) error {
	return s.repo.Add(ctx, team)
}

func (s *teamService) Get(teamName string) ([]dto.TeamMember, error) {
	return s.repo.Get(teamName)
}

func (s *teamService) DeactivateMembers(ctx context.Context, req dto.TeamDeactivateRequest) (*dto.TeamDeactivateResponse, error) {
	for _, userID := range req.UserIDs {
		prIDs, err := s.prRepo.ListOpenAssignments(ctx, userID)
		if err != nil {
			return nil, err
		}

		for _, prID := range prIDs {
			_, _, err := s.prRepo.Reassign(ctx, dto.ReassignRequest{
				PullRequestID: prID,
				OldUserID:     userID,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	if err := s.repo.DeactivateMembers(ctx, req.TeamName, req.UserIDs); err != nil {
		return nil, err
	}

	return &dto.TeamDeactivateResponse{
		TeamName:    req.TeamName,
		Deactivated: req.UserIDs,
	}, nil
}
