package repository

import "pr-reviwer-assigner/internal/domain/dto"

type TeamRepository interface {
	Add(dto.Team) error
	Get() (dto.Team, error)
}
