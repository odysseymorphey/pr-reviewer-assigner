package handlers_test

import (
	"bytes"
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

type userServiceMock struct {
	getReviewFn func(userID string) ([]dto.PRShort, error)
	setFn       func(req dto.SIARequest) (*dto.User, error)
}

func (m *userServiceMock) GetReview(userID string) ([]dto.PRShort, error) {
	if m.getReviewFn == nil {
		return nil, nil
	}
	return m.getReviewFn(userID)
}

func (m *userServiceMock) SetIsActive(req dto.SIARequest) (*dto.User, error) {
	if m.setFn == nil {
		return nil, nil
	}
	return m.setFn(req)
}

func TestUserHandlerSetIsActive_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &userServiceMock{
		setFn: func(req dto.SIARequest) (*dto.User, error) {
			return &dto.User{
				ID:       req.ID,
				Name:     "Alice",
				Team:     "backend",
				IsActive: req.IsActive,
			}, nil
		},
	}
	h := handlers.NewUserHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/users/setIsActive", h.SetIsActive)

	payload := []byte(`{"user_id":"u1","is_active":true}`)
	req := httptest.NewRequest("POST", "/users/setIsActive", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body dto.UserResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "u1", body.User.ID)
}

func TestUserHandlerSetIsActive_NotFound(t *testing.T) {
	app := fiber.New()
	mockSvc := &userServiceMock{
		setFn: func(req dto.SIARequest) (*dto.User, error) {
			return nil, errors2.ErrNotFound
		},
	}
	h := handlers.NewUserHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/users/setIsActive", h.SetIsActive)

	payload := []byte(`{"user_id":"ghost","is_active":false}`)
	req := httptest.NewRequest("POST", "/users/setIsActive", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestUserHandlerGetReview_BadRequest(t *testing.T) {
	app := fiber.New()
	h := handlers.NewUserHandler(&userServiceMock{}, zap.NewNop().Sugar())
	app.Get("/users/getReview", h.GetReview)

	req := httptest.NewRequest("GET", "/users/getReview", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserHandlerGetReview_NotFound(t *testing.T) {
	app := fiber.New()
	mockSvc := &userServiceMock{
		getReviewFn: func(userID string) ([]dto.PRShort, error) {
			return nil, errors2.ErrNotFound
		},
	}
	h := handlers.NewUserHandler(mockSvc, zap.NewNop().Sugar())
	app.Get("/users/getReview", h.GetReview)

	req := httptest.NewRequest("GET", "/users/getReview?user_id=ghost", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestUserHandlerGetReview_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &userServiceMock{
		getReviewFn: func(userID string) ([]dto.PRShort, error) {
			return []dto.PRShort{
				{ID: "pr-1", Name: "Title", AuthorID: "u1", Status: "OPEN"},
			}, nil
		},
	}
	h := handlers.NewUserHandler(mockSvc, zap.NewNop().Sugar())
	app.Get("/users/getReview", h.GetReview)

	req := httptest.NewRequest("GET", "/users/getReview?user_id=u2", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var out dto.UserPR
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	require.Equal(t, "u2", out.ID)
	require.Len(t, out.PRs, 1)
}
