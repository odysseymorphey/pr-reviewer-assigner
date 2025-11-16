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

type teamServiceMock struct {
	addFn        func(ctx context.Context, team dto.Team) error
	getFn        func(teamName string) ([]dto.TeamMember, error)
	deactivateFn func(ctx context.Context, req dto.TeamDeactivateRequest) (*dto.TeamDeactivateResponse, error)
}

func (m *teamServiceMock) Add(ctx context.Context, team dto.Team) error {
	if m.addFn == nil {
		return nil
	}
	return m.addFn(ctx, team)
}

func (m *teamServiceMock) Get(teamName string) ([]dto.TeamMember, error) {
	if m.getFn == nil {
		return nil, nil
	}
	return m.getFn(teamName)
}

func (m *teamServiceMock) DeactivateMembers(ctx context.Context, req dto.TeamDeactivateRequest) (*dto.TeamDeactivateResponse, error) {
	if m.deactivateFn == nil {
		return &dto.TeamDeactivateResponse{}, nil
	}
	return m.deactivateFn(ctx, req)
}

func TestTeamHandlerGet_BadRequest(t *testing.T) {
	app := fiber.New()
	h := handlers.NewTeamHandler(&teamServiceMock{}, zap.NewNop().Sugar())
	app.Get("/team/get", h.Get)

	req := httptest.NewRequest("GET", "/team/get", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTeamHandlerGet_NotFound(t *testing.T) {
	app := fiber.New()
	mockSvc := &teamServiceMock{
		getFn: func(teamName string) ([]dto.TeamMember, error) {
			return nil, errors2.ErrNotFound
		},
	}
	h := handlers.NewTeamHandler(mockSvc, zap.NewNop().Sugar())
	app.Get("/team/get", h.Get)

	req := httptest.NewRequest("GET", "/team/get?team_name=ghosts", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestTeamHandlerGet_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &teamServiceMock{
		getFn: func(teamName string) ([]dto.TeamMember, error) {
			return []dto.TeamMember{
				{ID: "u1", Name: "Alice", IsActive: true},
			}, nil
		},
	}
	h := handlers.NewTeamHandler(mockSvc, zap.NewNop().Sugar())
	app.Get("/team/get", h.Get)

	req := httptest.NewRequest("GET", "/team/get?team_name=backend", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body dto.Team
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, "backend", body.Name)
	require.Len(t, body.Members, 1)
}

func TestTeamHandlerAdd_ValidationError(t *testing.T) {
	app := fiber.New()
	h := handlers.NewTeamHandler(&teamServiceMock{}, zap.NewNop().Sugar())
	app.Post("/team/add", h.Add)

	payload := []byte(`{"team_name":"","members":[]}`)
	req := httptest.NewRequest("POST", "/team/add", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTeamHandlerAdd_Success(t *testing.T) {
	app := fiber.New()
	addCalled := false
	mockSvc := &teamServiceMock{
		addFn: func(ctx context.Context, team dto.Team) error {
			addCalled = true
			require.Equal(t, "backend", team.Name)
			require.Len(t, team.Members, 1)
			return nil
		},
	}
	h := handlers.NewTeamHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/team/add", h.Add)

	body := dto.Team{
		Name: "backend",
		Members: []dto.TeamMember{
			{ID: "u1", Name: "Alice", IsActive: true},
		},
	}
	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/team/add", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.True(t, addCalled)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestTeamHandlerDeactivateMembers_Validation(t *testing.T) {
	app := fiber.New()
	h := handlers.NewTeamHandler(&teamServiceMock{}, zap.NewNop().Sugar())
	app.Post("/team/deactivateMembers", h.DeactivateMembers)

	req := httptest.NewRequest("POST", "/team/deactivateMembers", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTeamHandlerDeactivateMembers_Success(t *testing.T) {
	app := fiber.New()
	mockSvc := &teamServiceMock{
		deactivateFn: func(ctx context.Context, req dto.TeamDeactivateRequest) (*dto.TeamDeactivateResponse, error) {
			require.Equal(t, "backend", req.TeamName)
			require.Equal(t, []string{"u1"}, req.UserIDs)
			return &dto.TeamDeactivateResponse{
				TeamName:    "backend",
				Deactivated: []string{"u1"},
			}, nil
		},
	}
	h := handlers.NewTeamHandler(mockSvc, zap.NewNop().Sugar())
	app.Post("/team/deactivateMembers", h.DeactivateMembers)

	body := []byte(`{"team_name":"backend","user_ids":["u1"]}`)
	req := httptest.NewRequest("POST", "/team/deactivateMembers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}
