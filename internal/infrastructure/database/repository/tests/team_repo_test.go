package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"pr-reviwer-assigner/internal/domain/dto"
	errors2 "pr-reviwer-assigner/internal/errors"
	repo "pr-reviwer-assigner/internal/infrastructure/database/repository"
)

func TestTeamRepoGet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewTeamRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "username", "is_active"}).
		AddRow("u1", "Alice", true).
		AddRow("u2", "Bob", false)

	mock.ExpectQuery(`SELECT user_id, username, is_active FROM users WHERE team_name = \$1`).
		WithArgs("backend").
		WillReturnRows(rows)

	members, err := r.Get("backend")
	require.NoError(t, err)
	require.Len(t, members, 2)
	require.Equal(t, "u1", members[0].ID)
	require.Equal(t, false, members[1].IsActive)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepoGet_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewTeamRepository(db)

	mock.ExpectQuery(`SELECT user_id, username, is_active FROM users WHERE team_name = \$1`).
		WithArgs("ghosts").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "is_active"}))

	mock.ExpectQuery(`SELECT 1 FROM teams WHERE team_name = \$1`).
		WithArgs("ghosts").
		WillReturnError(sql.ErrNoRows)

	_, err = r.Get("ghosts")
	require.ErrorIs(t, err, errors2.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepoAdd_AlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewTeamRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO teams \(team_name\) VALUES \(\$1\)`).
		WithArgs("backend").
		WillReturnError(&pq.Error{Code: "23505"})
	mock.ExpectRollback()

	team := dto.Team{
		Name: "backend",
		Members: []dto.TeamMember{
			{ID: "u1", Name: "Alice", IsActive: true},
		},
	}

	err = r.Add(context.Background(), team)
	require.ErrorIs(t, err, errors2.ErrTeamExists)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTeamRepoAdd_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewTeamRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO teams \(team_name\) VALUES \(\$1\)`).
		WithArgs("backend").
		WillReturnResult(sqlmock.NewResult(0, 1))

	userUpsert := regexp.QuoteMeta(`
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
   			SET username = EXCLUDED.username,
       			team_name = EXCLUDED.team_name,
       			is_active = EXCLUDED.is_active`)

	mock.ExpectExec(userUpsert).
		WithArgs("u1", "Alice", "backend", true).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(userUpsert).
		WithArgs("u2", "Bob", "backend", false).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	team := dto.Team{
		Name: "backend",
		Members: []dto.TeamMember{
			{ID: "u1", Name: "Alice", IsActive: true},
			{ID: "u2", Name: "Bob", IsActive: false},
		},
	}

	err = r.Add(context.Background(), team)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
