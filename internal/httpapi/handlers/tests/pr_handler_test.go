package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"pr-reviwer-assigner/internal/httpapi/handlers"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"pr-reviwer-assigner/internal/domain/dto"
	errors2 "pr-reviwer-assigner/internal/errors"
)

type prServiceMock struct {
	createFn   func(ctx context.Context, req dto.PRRequest) (dto.PR, error)
	mergeFn    func(ctx context.Context, req dto.MergeRequest) (*dto.PR, error)
	reassignFn func(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error)
}

func (m *prServiceMock) Create(ctx context.Context, req dto.PRRequest) (dto.PR, error) {
	if m.createFn == nil {
		return dto.PR{}, nil
	}
	return m.createFn(ctx, req)
}

func (m *prServiceMock) Merge(ctx context.Context, req dto.MergeRequest) (*dto.PR, error) {
	if m.mergeFn == nil {
		return nil, nil
	}
	return m.mergeFn(ctx, req)
}

func (m *prServiceMock) Reassign(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error) {
	if m.reassignFn == nil {
		return nil, "", nil
	}
	return m.reassignFn(ctx, req)
}

func TestPRHandlerCreate_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &prServiceMock{
		createFn: func(ctx context.Context, req dto.PRRequest) (dto.PR, error) {
			return dto.PR{
				ID:        req.ID,
				Name:      req.Name,
				AuthorID:  req.AuthorID,
				Status:    "OPEN",
				Reviewers: []string{"u2"},
			}, nil
		},
	}
	h := handlers.NewPRHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/pullRequest/create", h.CreatePR)

	payload := []byte(`{"pull_request_id":"pr-1","pull_request_name":"Add","author_id":"u1"}`)
	req := httptest.NewRequest("POST", "/pullRequest/create", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var out dto.PRResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	require.Equal(t, "pr-1", out.PR.ID)
	require.Len(t, out.PR.Reviewers, 1)
}

func TestPRHandlerCreate_NotFound(t *testing.T) {
	app := fiber.New()
	mockSvc := &prServiceMock{
		createFn: func(ctx context.Context, req dto.PRRequest) (dto.PR, error) {
			return dto.PR{}, errors2.ErrNotFound
		},
	}
	h := handlers.NewPRHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/pullRequest/create", h.CreatePR)

	payload := []byte(`{"pull_request_id":"pr-1","pull_request_name":"Add","author_id":"u1"}`)
	req := httptest.NewRequest("POST", "/pullRequest/create", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestPRHandlerMerge_BadRequest(t *testing.T) {
	app := fiber.New()
	h := handlers.NewPRHandler(&prServiceMock{}, zap.NewNop().Sugar())
	app.Post("/pullRequest/merge", h.MergePR)

	req := httptest.NewRequest("POST", "/pullRequest/merge", bytes.NewReader([]byte(`{"pull_request_id":""}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestPRHandlerMerge_NotFound(t *testing.T) {
	app := fiber.New()
	mockSvc := &prServiceMock{
		mergeFn: func(ctx context.Context, req dto.MergeRequest) (*dto.PR, error) {
			return nil, errors2.ErrNotFound
		},
	}
	h := handlers.NewPRHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/pullRequest/merge", h.MergePR)

	req := httptest.NewRequest("POST", "/pullRequest/merge", bytes.NewReader([]byte(`{"pull_request_id":"unknown"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestPRHandlerReassign_Conflict(t *testing.T) {
	app := fiber.New()
	mockSvc := &prServiceMock{
		reassignFn: func(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error) {
			return nil, "", errors2.ErrPRMerged
		},
	}
	h := handlers.NewPRHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/pullRequest/reassign", h.ReassignViewer)

	req := httptest.NewRequest("POST", "/pullRequest/reassign", bytes.NewReader([]byte(`{"pull_request_id":"pr-1","old_user_id":"u2"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestPRHandlerReassign_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &prServiceMock{
		reassignFn: func(ctx context.Context, req dto.ReassignRequest) (*dto.PR, string, error) {
			return &dto.PR{
				ID:        req.PullRequestID,
				Name:      "PR",
				AuthorID:  "u1",
				Status:    "OPEN",
				Reviewers: []string{"u3"},
			}, "u3", nil
		},
	}
	h := handlers.NewPRHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/pullRequest/reassign", h.ReassignViewer)

	req := httptest.NewRequest("POST", "/pullRequest/reassign", bytes.NewReader([]byte(`{"pull_request_id":"pr-1","old_user_id":"u2"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body dto.ReassignResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "u3", body.ReplacedBy)
}
