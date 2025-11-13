package repository

import "pr-reviwer-assigner/internal/domain/dto"

type TeamRepository interface {
	Add(team dto.Team) error
	Get(teamName string) (dto.Team, error)
}
