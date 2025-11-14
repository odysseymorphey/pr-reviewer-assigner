package handlers

import (
	"encoding/json"
	"errors"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/services"
	errors2 "pr-reviwer-assigner/internal/errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type PRHandler struct {
	service services.PRService
	logger  *zap.SugaredLogger
}

func NewPRHandler(service services.PRService, logger *zap.SugaredLogger) *PRHandler {
	return &PRHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PRHandler) CreatePR(c fiber.Ctx) error {
	var prReq dto.PRRequest

	err := json.Unmarshal(c.Body(), &prReq)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
	}

	prReq.ID = strings.TrimSpace(prReq.ID)
	prReq.Name = strings.TrimSpace(prReq.Name)
	prReq.AuthorID = strings.TrimSpace(prReq.AuthorID)

	if prReq.ID == "" || prReq.Name == "" || prReq.AuthorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "all fields are required",
			},
		})
	}

	ctx := c.Context()

	pr, err := h.service.Create(ctx, prReq)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrPRExists):
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "PR already exists",
				},
			})
		case errors.Is(err, errors2.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
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

	response := dto.PRResponse{
		PR: pr,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *PRHandler) MergePR(c fiber.Ctx) error {
	var req dto.MergeRequest

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
	}

	req.PullRequestID = strings.TrimSpace(req.PullRequestID)
	if req.PullRequestID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "missing pull request id",
			},
		})
	}

	ctx := c.Context()

	pr, err := h.service.Merge(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
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

	response := dto.PRResponse{
		PR: *pr,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PRHandler) ReassignViewer(c fiber.Ctx) error {
	var req dto.ReassignRequest

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
	}

	req.PullRequestID = strings.TrimSpace(req.PullRequestID)
	req.OldUserID = strings.TrimSpace(req.OldUserID)
	if req.PullRequestID == "" || req.OldUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "all fields are required",
			},
		})
	}

	ctx := c.Context()

	pr, replacedBy, err := h.service.Reassign(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "resource not found",
				},
			})
		case errors.Is(err, errors2.ErrPRMerged):
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "cannot reassign on merged PR",
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

	response := dto.ReassignResponse{
		PR: dto.PR{
			ID:        pr.ID,
			Name:      pr.Name,
			AuthorID:  pr.AuthorID,
			Status:    pr.Status,
			Reviewers: pr.Reviewers,
		},
		ReplacedBy: replacedBy,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
