package services

import (
	"context"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type PRService interface {
	Create(ctx context.Context, req dto.PRRequest) (dto.PR, error)
	Merge(ctx context.Context, req dto.MergeRequest) (*dto.PR, error)
	Reassign(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error)
}

type prService struct {
	repo repository.PRRepository
}

func NewPRService(repo repository.PRRepository) PRService {
	return &prService{
		repo: repo,
	}
}

func (s *prService) Create(ctx context.Context, req dto.PRRequest) (dto.PR, error) {
	users, err := s.repo.Create(ctx, req)
	if err != nil {
		return dto.PR{}, err
	}

	var pr dto.PR
	pr.ID = req.ID
	pr.Name = req.Name
	pr.AuthorID = req.AuthorID
	pr.Status = "OPEN"
	pr.Reviewers = make([]string, 0)
	for _, user := range users {
		pr.Reviewers = append(pr.Reviewers, user)
	}

	return pr, nil
}

func (s *prService) Merge(ctx context.Context, req dto.MergeRequest) (*dto.PR, error) {
	return s.repo.Merge(ctx, req)
}

func (s *prService) Reassign(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error) {
	return s.repo.Reassign(ctx, req)
}
