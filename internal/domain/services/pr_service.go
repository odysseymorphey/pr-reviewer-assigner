package services

import (
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type PRService interface {
	Create(in dto.PRRequest)
	Merge()
	Reassign()
}

type prService struct {
	repo repository.PRRepository
}

func NewPRService(repo repository.PRRepository) PRService {
	return &prService{
		repo: repo,
	}
}

func (s *prService) Create(in dto.PRRequest) {

}

func (s *prService) Merge() {

}

func (s *prService) Reassign() {

}
