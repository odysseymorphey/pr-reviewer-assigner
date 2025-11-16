package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"pr-reviwer-assigner/internal/domain/dto"
	errors2 "pr-reviwer-assigner/internal/errors"
	repo "pr-reviwer-assigner/internal/infrastructure/database/repository"
)

func TestPRRepoReassign_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewPRRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT\s+pull_request_id`).
		WithArgs("pr-1").
		WillReturnRows(sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
			AddRow("pr-1", "Add search", "author-1", "OPEN"))

	mock.ExpectQuery(`SELECT\s+team_name\s+FROM users WHERE user_id = \$1`).
		WithArgs("old-user").
		WillReturnRows(sqlmock.NewRows([]string{"team_name"}).AddRow("backend"))

	mock.ExpectQuery(`SELECT 1\s+FROM pull_request_reviewers`).
		WithArgs("pr-1", "old-user").
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	mock.ExpectQuery(`SELECT u\.user_id\s+FROM users u`).
		WithArgs("backend", "author-1", "old-user", "pr-1").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("new-user"))

	mock.ExpectExec(`DELETE FROM pull_request_reviewers`).
		WithArgs("pr-1", "old-user").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`INSERT INTO pull_request_reviewers`).
		WithArgs("pr-1", "new-user").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(`SELECT reviewer_id FROM pull_request_reviewers`).
		WithArgs("pr-1").
		WillReturnRows(sqlmock.NewRows([]string{"reviewer_id"}).
			AddRow("new-user").
			AddRow("another"))

	mock.ExpectCommit()

	req := dto.ReassignRequest{
		PullRequestID: "pr-1",
		OldUserID:     "old-user",
	}

	pr, newReviewer, err := r.Reassign(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "new-user", newReviewer)
	require.ElementsMatch(t, []string{"another", "new-user"}, pr.Reviewers)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPRRepoReassign_PRMerged(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := repo.NewPRRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT\s+pull_request_id`).
		WithArgs("pr-merged").
		WillReturnRows(sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
			AddRow("pr-merged", "Title", "author-1", "MERGED"))
	mock.ExpectRollback()

	req := dto.ReassignRequest{
		PullRequestID: "pr-merged",
		OldUserID:     "old-user",
	}

	_, _, err = r.Reassign(context.Background(), req)
	require.ErrorIs(t, err, errors2.ErrPRMerged)
	require.NoError(t, mock.ExpectationsWereMet())
}
