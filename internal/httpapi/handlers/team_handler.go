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

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		h.logger.Error("team add: failed to unmarshal body: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrInternal.Error(),
				Message: "internal server error",
			},
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		h.logger.Error("team add: empty team_name")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "name can't be empty",
			},
		})
	}

	if len(req.Members) == 0 {
		h.logger.Error("team add: members empty")
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
			h.logger.Error("team add: member empty fields: ", i)
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrBadRequest.Error(),
					Message: fmt.Sprintf("member[%d]: user_id and username can't be empty", i),
				},
			})
		}

		if _, ok := seen[m.ID]; ok {
			h.logger.Error("team add: duplicate member: ", m.ID)
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
			h.logger.Error("team add: team exists: ", req.Name)
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrTeamExists.Error(),
					Message: "team_name already exists",
				},
			})
		default:
			h.logger.Error("team add: service error: ", err)
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
	h.logger.Info("team add success: ", req.Name)

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *TeamHandler) Get(c fiber.Ctx) error {
	teamName := strings.TrimSpace(c.Query("team_name"))
	if teamName == "" {
		h.logger.Error("team get: empty team_name")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "team_name can't be empty",
			},
		})
	}

	members, err := h.teamService.Get(teamName)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrNotFound):
			h.logger.Error("team get: not found: ", teamName)
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrNotFound.Error(),
					Message: "resource not found",
				},
			})
		default:
			h.logger.Error("team get: service error: ", err)
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
	h.logger.Info("team get success: ", teamName)

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *TeamHandler) DeactivateMembers(c fiber.Ctx) error {
	var req dto.TeamDeactivateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		h.logger.Error("team deactivate: failed to unmarshal body: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrInternal.Error(),
				Message: "internal server error",
			},
		})
	}

	req.TeamName = strings.TrimSpace(req.TeamName)
	if req.TeamName == "" {
		h.logger.Error("team deactivate: empty team_name")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "team_name can't be empty",
			},
		})
	}

	if len(req.UserIDs) == 0 {
		h.logger.Error("team deactivate: empty user_ids")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "user_ids can't be empty",
			},
		})
	}

	seen := make(map[string]struct{})
	for i, id := range req.UserIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			h.logger.Error("team deactivate: empty user_id at index: ", i)
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrBadRequest.Error(),
					Message: "user_ids can't contain empty values",
				},
			})
		}
		if _, ok := seen[id]; ok {
			h.logger.Error("team deactivate: duplicate user_id: ", id)
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrBadRequest.Error(),
					Message: "user_ids must be unique",
				},
			})
		}
		seen[id] = struct{}{}
		req.UserIDs[i] = id
	}

	resp, err := h.teamService.DeactivateMembers(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrNotFound):
			h.logger.Error("team deactivate: not found: ", err)
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrNotFound.Error(),
					Message: "resource not found",
				},
			})
		case errors.Is(err, errors2.ErrNoCandidate):
			h.logger.Error("team deactivate: no candidate: ", err)
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrNoCandidate.Error(),
					Message: "no active replacement candidate in team",
				},
			})
		default:
			h.logger.Error("team deactivate: internal error: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrInternal.Error(),
					Message: "internal server error",
				},
			})
		}
	}

	h.logger.Info("team deactivate success: ", resp)

	return c.Status(fiber.StatusOK).JSON(resp)
}
