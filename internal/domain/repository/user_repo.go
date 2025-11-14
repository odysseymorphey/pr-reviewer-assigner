package repository

import "pr-reviwer-assigner/internal/domain/dto"

type UserRepository interface {
	GetReview(userID string) ([]dto.PRShort, error)
	SetIsActive(user dto.SIARequest) (*dto.User, error)
}
