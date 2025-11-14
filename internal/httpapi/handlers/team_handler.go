package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/services"
	errors2 "pr-reviwer-assigner/internal/errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type TeamHandler struct {
	teamService services.TeamService
	logger      *zap.SugaredLogger
}

func NewTeamHandler(teamService services.TeamService, logger *zap.SugaredLogger) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		logger:      logger,
	}
}

func (h *TeamHandler) Add(c fiber.Ctx) error {
	var req dto.Team

	_ = json.Unmarshal(c.Body(), &req)

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "name can't be empty",
			},
		})
	}

	if len(req.Members) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "members can't be empty",
			},
		})
	}

	seen := make(map[string]struct{})
	for i, m := range req.Members {
		m.ID = strings.TrimSpace(m.ID)
		m.Name = strings.TrimSpace(m.Name)

		if m.ID == "" || m.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrBadRequest.Error(),
					Message: fmt.Sprintf("member[%d]: user_id and username can't be empty", i),
				},
			})
		}

		if _, ok := seen[m.ID]; ok {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrBadRequest.Error(),
					Message: fmt.Sprintf("duplicate user_id in members: %s", m.ID),
				},
			})
		}
		seen[m.ID] = struct{}{}

		req.Members[i] = m
	}

	ctx := c.Context()

	err := h.teamService.Add(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrTeamExists):
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrTeamExists.Error(),
					Message: "team_name already exists",
				},
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrInternal.Error(),
					Message: "internal server error",
				},
			})
		}
	}

	response := dto.TeamResponse{
		Team: req,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *TeamHandler) Get(c fiber.Ctx) error {
	teamName := strings.TrimSpace(c.Params("team_name"))
	if teamName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	members, err := h.teamService.Get(teamName)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrNotFound.Error(),
					Message: "resource not found",
				},
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrInternal.Error(),
					Message: "internal server error",
				},
			})
		}
	}

	response := &dto.Team{
		Name:    teamName,
		Members: members,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
