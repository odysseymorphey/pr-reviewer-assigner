package repository

import (
	"database/sql"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (s *userRepo) GetReview(userID string) (dto.UserPR, error) {
	return dto.UserPR{}, nil
}

func (s *userRepo) SetIsActive(user dto.User) (dto.UserResponse, error) {
	return dto.UserResponse{}, nil
}
