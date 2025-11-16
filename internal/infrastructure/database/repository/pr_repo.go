package repository

import (
	"context"
	"database/sql"
	"errors"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/repository"
	errors2 "pr-reviwer-assigner/internal/errors"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type prRepo struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) repository.PRRepository {
	return &prRepo{
		db: db,
	}
}

func (s *prRepo) Create(ctx context.Context, req dto.PRRequest) ([]string, error) {
	const createQuery = `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	const pickupQuery = `
		SELECT user_id
		FROM users
		WHERE team_name = $1
			AND is_active = TRUE
			AND user_id <> $2
		ORDER BY user_id
		LIMIT 2
	`

	const insertReviewerQuery = `
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var authorTeam string
	var authorActive bool
	err = tx.QueryRowContext(ctx,
		`SELECT team_name, is_active
		   FROM users
		  WHERE user_id = $1`,
		req.AuthorID,
	).Scan(&authorTeam, &authorActive)

	if errors.Is(err, sql.ErrNoRows) || !authorActive {
		return nil, errors2.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UTC().String()
	_, err = tx.ExecContext(ctx, createQuery,
		req.ID,
		req.Name,
		req.AuthorID,
		"OPEN",
		createdAt,
	)
	if err != nil {
		switch {
		case isUniqueViolation(err):
			return nil, errors2.ErrPRExists
		default:
			return nil, err
		}
	}

	rows, err := tx.QueryContext(ctx, pickupQuery,
		authorTeam,
		req.AuthorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, id := range reviewers {
		if _, err := tx.ExecContext(ctx, insertReviewerQuery, req.ID, id); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return reviewers, nil
}

func (s *prRepo) Merge(ctx context.Context, req dto.MergeRequest) (*dto.PR, error) {
	const mergeQuery = `
		UPDATE pull_requests
   		SET status    = $2,
       		merged_at = COALESCE(merged_at, $3)
 		WHERE pull_request_id = $1
		RETURNING 
		    pull_request_id,
            pull_request_name,
            author_id,
            status::text,
            created_at,
            merged_at
 		`

	const reviewersQuery = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
		ORDER BY reviewer_id
		`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	mergedAt := time.Now().UTC().String()
	var pr dto.PR
	err = tx.QueryRowContext(ctx, mergeQuery, req.PullRequestID, "MERGED", mergedAt).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors2.ErrNotFound
		default:
			return nil, err
		}
	}

	rows, err := tx.QueryContext(ctx, reviewersQuery, pr.ID)
	if err != nil {
		return nil, err
	}

	var reviewers []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	pr.Reviewers = reviewers

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (s *prRepo) Reassign(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error) {
	const prQuery = `
		SELECT 
		    pull_request_id,
			pull_request_name,
			author_id,
			status::text
		FROM pull_requests
		WHERE pull_request_id = $1
		FOR UPDATE
	`

	const oldUserQuery = `
		SELECT 
		    team_name
		FROM users
		WHERE user_id = $1
	`

	const checkOldUserQuery = `
		SELECT 1
	  	FROM pull_request_reviewers
		WHERE 
		    pull_request_id = $1
			AND reviewer_id = $2
	`

	const findNewRevQuery = `
		SELECT u.user_id
		FROM users u
		WHERE u.team_name = $1
			AND u.is_active = TRUE
			AND u.user_id <> $2
		    AND u.user_id <> $3
		    AND NOT EXISTS (
			   	SELECT 1
					FROM pull_request_reviewers prr
				WHERE prr.pull_request_id = $4
					AND prr.reviewer_id     = u.user_id
		   )
		ORDER BY u.user_id
		LIMIT 1
	`

	const deleteOldRevQuery = `
		DELETE FROM pull_request_reviewers
		WHERE pull_request_id = $1
		AND reviewer_id     = $2
	`

	const insertNewRevQuery = `
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	const collectRevQuery = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
		ORDER BY reviewer_id
	`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	var pr dto.PR
	err = tx.QueryRowContext(ctx, prQuery, req.PullRequestID).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, "", errors2.ErrNotFound
		default:
			return nil, "", err
		}
	}

	if pr.Status == "MERGED" {
		return nil, "", errors2.ErrPRMerged
	}

	var oldUserTeam string
	err = tx.QueryRowContext(ctx, oldUserQuery, req.OldUserID).Scan(&oldUserTeam)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, "", errors2.ErrNotFound
		default:
			return nil, "", err
		}
	}

	var dummy int
	err = tx.QueryRowContext(ctx, checkOldUserQuery, req.PullRequestID, req.OldUserID).Scan(&dummy)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, "", errors2.ErrNotAssigned
		default:
			return nil, "", err
		}
	}

	var newReviewerID string
	err = tx.QueryRowContext(ctx, findNewRevQuery, oldUserTeam, pr.AuthorID, req.OldUserID, req.PullRequestID).Scan(&newReviewerID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, "", errors2.ErrNoCandidate
		default:
			return nil, "", err
		}
	}

	_, err = tx.ExecContext(ctx, deleteOldRevQuery, req.PullRequestID, req.OldUserID)
	if err != nil {
		return nil, "", err
	}

	_, err = tx.ExecContext(ctx, insertNewRevQuery, req.PullRequestID, newReviewerID)
	if err != nil {
		return nil, "", err
	}

	rows, err := tx.QueryContext(ctx, collectRevQuery, req.PullRequestID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, "", err
		}
		reviewers = append(reviewers, reviewerID)
	}

	pr.Reviewers = reviewers

	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	if err := tx.Commit(); err != nil {
		return nil, "", err
	}

	return &pr, newReviewerID, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}

func (s *prRepo) ListOpenAssignments(ctx context.Context, reviewerID string) ([]string, error) {
	const query = `
		SELECT pr.pull_request_id
		FROM pull_requests pr
		JOIN pull_request_reviewers prr ON prr.pull_request_id = pr.pull_request_id
		WHERE pr.status = 'OPEN'
		  AND prr.reviewer_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
