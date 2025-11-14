package repository

import (
	"context"
	"pr-reviwer-assigner/internal/domain/dto"
)

type PRRepository interface {
	Create(ctx context.Context, req dto.PRRequest) ([]string, error)
	Merge(ctx context.Context, req dto.MergeRequest) (*dto.PR, error)
	Reassign(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error)
}
