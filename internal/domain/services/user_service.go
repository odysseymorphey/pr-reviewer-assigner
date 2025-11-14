package services

import (
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type UserService interface {
	GetReview(userID string) ([]dto.PRShort, error)
	SetIsActive(req dto.SIARequest) (*dto.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetReview(userID string) ([]dto.PRShort, error) {
	return s.repo.GetReview(userID)
}

func (s *userService) SetIsActive(user dto.SIARequest) (*dto.User, error) {
	u, err := s.repo.SetIsActive(user)
	if err != nil {
		return nil, err
	}

	return u, nil
}
