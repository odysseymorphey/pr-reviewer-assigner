package repository

import (
	"database/sql"
	"errors"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
	errors2 "pr-reviwer-assigner/internal/errors"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (s *userRepo) GetReview(userID string) ([]dto.PRShort, error) {
	const getReview = `
		SELECT
			pr.pull_request_id,
			pr.pull_request_name,
			pr.author_id,
			pr.status::text
		FROM pull_requests pr
		JOIN pull_request_reviewers prr
		  ON prr.pull_request_id = pr.pull_request_id
		WHERE prr.reviewer_id = $1
		ORDER BY pr.created_at NULLS LAST, pr.pull_request_id`

	rows, err := s.db.Query(getReview, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []dto.PRShort
	for rows.Next() {
		var pr dto.PRShort
		err = rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
		)

		if err != nil {
			return nil, err
		}

		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		const userExists = `SELECT 1 FROM users WHERE user_id = $1`
		var dummy int
		err = s.db.QueryRow(userExists, userID).Scan(&dummy)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, errors2.ErrNotFound
			default:
				return nil, err
			}
		}
	}

	return prs, nil
}

func (s *userRepo) SetIsActive(req dto.SIARequest) (*dto.User, error) {
	const setIsActive = `
		UPDATE users
		   SET is_active = $1
		 WHERE user_id   = $2
		RETURNING user_id, username, team_name, is_active
	`
	var user dto.User
	err := s.db.QueryRow(setIsActive, req.IsActive, req.ID).Scan(
		&user.ID,
		&user.Name,
		&user.Team,
		&user.IsActive,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors2.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
