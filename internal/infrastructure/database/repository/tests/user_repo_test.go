package repository_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"pr-reviwer-assigner/internal/domain/dto"
	errors2 "pr-reviwer-assigner/internal/errors"
	repo "pr-reviwer-assigner/internal/infrastructure/database/repository"
)

func TestUserRepoSetIsActive_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewUserRepository(db)

	row := sqlmock.NewRows([]string{"user_id", "username", "team_name", "is_active"}).
		AddRow("u1", "Alice", "backend", true)

	mock.ExpectQuery(`UPDATE users\s+SET is_active = \$1\s+WHERE user_id = \$2`).
		WithArgs(true, "u1").
		WillReturnRows(row)

	req := dto.SIARequest{
		ID:       "u1",
		IsActive: true,
	}

	user, err := r.SetIsActive(req)
	require.NoError(t, err)
	require.Equal(t, "backend", user.Team)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepoSetIsActive_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewUserRepository(db)

	mock.ExpectQuery(`UPDATE users\s+SET is_active = \$1\s+WHERE user_id = \$2`).
		WithArgs(false, "ghost").
		WillReturnError(sql.ErrNoRows)

	req := dto.SIARequest{
		ID:       "ghost",
		IsActive: false,
	}

	_, err = r.SetIsActive(req)
	require.ErrorIs(t, err, errors2.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepoGetReview_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewUserRepository(db)

	mock.ExpectQuery(`SELECT\s+pr\.pull_request_id`).
		WithArgs("ghost").
		WillReturnRows(sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}))

	mock.ExpectQuery(`SELECT 1 FROM users WHERE user_id = \$1`).
		WithArgs("ghost").
		WillReturnError(sql.ErrNoRows)

	_, err = r.GetReview("ghost")
	require.ErrorIs(t, err, errors2.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepoGetReview_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
		AddRow("pr-1", "Add login", "u1", "OPEN").
		AddRow("pr-2", "Fix tests", "u1", "MERGED")

	mock.ExpectQuery(`SELECT\s+pr\.pull_request_id`).
		WithArgs("u2").
		WillReturnRows(rows)

	prs, err := r.GetReview("u2")
	require.NoError(t, err)
	require.Len(t, prs, 2)
	require.Equal(t, "pr-2", prs[1].ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
