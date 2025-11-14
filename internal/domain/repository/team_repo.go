package repository

import (
	"context"
	"pr-reviwer-assigner/internal/domain/dto"
)

type TeamRepository interface {
	Add(ctx context.Context, team dto.Team) error
	Get(teamName string) ([]dto.TeamMember, error)
}
