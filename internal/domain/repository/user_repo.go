package repository

import "pr-reviwer-assigner/internal/domain/dto"

type UserRepository interface {
	GetReview(userID string) (dto.UserPR, error)
	SetIsActive(user dto.User) (dto.UserResponse, error)
}
