package repository

import (
	"context"
	"database/sql"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
)

type teamRepo struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) repository.TeamRepository {
	return &teamRepo{
		db: db,
	}
}

func (r *teamRepo) Get(teamName string) ([]dto.TeamMember, error) {
	const query = `SELECT user_id, username, is_active FROM users WHERE team_name = $1`

	rows, err := r.db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []dto.TeamMember
	for rows.Next() {
		var member dto.TeamMember
		err = rows.Scan(
			&member.ID,
			&member.Name,
			&member.IsActive,
		)

		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (r *teamRepo) Add(ctx context.Context, team dto.Team) error {
	const query = `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
   			SET username = EXCLUDED.username,
       			team_name = EXCLUDED.team_name,
       			is_active = EXCLUDED.is_active`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO teams (team_name) VALUES ($1)`,
		team.Name,
	)

	for _, m := range team.Members {
		if _, err := tx.ExecContext(ctx,
			query,
			m.ID,
			m.Name,
			team.Name,
			m.IsActive,
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
