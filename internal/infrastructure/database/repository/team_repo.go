package repository

import (
	"context"
	"database/sql"
	"errors"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
	errors2 "pr-reviwer-assigner/internal/errors"
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(members) == 0 {
		const existsQuery = `SELECT 1 FROM teams WHERE team_name = $1`
		var dummy int
		err = r.db.QueryRow(existsQuery, teamName).Scan(&dummy)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, errors2.ErrNotFound
			default:
				return nil, err
			}
		}
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
	if err != nil {
		switch {
		case isUniqueViolation(err):
			return errors2.ErrTeamExists
		default:
			return err
		}
	}

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

func (r *teamRepo) DeactivateMembers(ctx context.Context, teamName string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	const query = `
		UPDATE users
		   SET is_active = FALSE
		 WHERE user_id = $1
		   AND team_name = $2
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, id := range userIDs {
		res, err := tx.ExecContext(ctx, query, id, teamName)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors2.ErrNotFound
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
