package services

import (
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type UserService interface {
	GetReview(userID string) (dto.UserPR, error)
	SetIsActive(user dto.User) (dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetReview(userID string) (dto.UserPR, error) {
	return s.repo.GetReview(userID)
}

func (s *userService) SetIsActive(user dto.User) (dto.UserResponse, error) {
	return s.repo.SetIsActive(user)
}
