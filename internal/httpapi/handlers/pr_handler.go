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

	if err := json.Unmarshal(c.Body(), &prReq); err != nil {
		h.logger.Error("create PR: failed to unmarshal body: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrInternal.Error(),
				Message: "internal server error",
			},
		})
	}

	prReq.ID = strings.TrimSpace(prReq.ID)
	prReq.Name = strings.TrimSpace(prReq.Name)
	prReq.AuthorID = strings.TrimSpace(prReq.AuthorID)

	if prReq.ID == "" || prReq.Name == "" || prReq.AuthorID == "" {
		h.logger.Error("create PR: missing fields: ", prReq)
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
			h.logger.Error("create PR: already exists: ", prReq.ID)
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "PR already exists",
				},
			})
		case errors.Is(err, errors2.ErrNotFound):
			h.logger.Error("create PR: not found author: ", prReq.AuthorID)
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "resource not found",
				},
			})
		default:
			h.logger.Error("create PR: service error: ", err)
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
	h.logger.Info("create PR success: ", pr.ID)

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *PRHandler) MergePR(c fiber.Ctx) error {
	var req dto.MergeRequest

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		h.logger.Error("merge PR: failed to unmarshal body: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrInternal.Error(),
				Message: "internal server error",
			},
		})
	}

	req.PullRequestID = strings.TrimSpace(req.PullRequestID)
	if req.PullRequestID == "" {
		h.logger.Error("merge PR: empty pull_request_id")
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
			h.logger.Error("merge PR: not found: ", req.PullRequestID)
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "resource not found",
				},
			})
		default:
			h.logger.Error("merge PR: service error: ", err)
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
	h.logger.Info("merge PR success: ", pr.ID)

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PRHandler) ReassignViewer(c fiber.Ctx) error {
	var req dto.ReassignRequest

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		h.logger.Error("reassign PR: failed to unmarshal body: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrInternal.Error(),
				Message: "internal server error",
			},
		})
	}

	req.PullRequestID = strings.TrimSpace(req.PullRequestID)
	req.OldUserID = strings.TrimSpace(req.OldUserID)
	if req.PullRequestID == "" || req.OldUserID == "" {
		h.logger.Error("reassign PR: missing fields: ", req)
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
			h.logger.Error("reassign PR: not found: ", req)
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "resource not found",
				},
			})
		case errors.Is(err, errors2.ErrPRMerged):
			h.logger.Error("reassign PR: merged: ", req.PullRequestID)
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    err.Error(),
					Message: "cannot reassign on merged PR",
				},
			})
		default:
			h.logger.Error("reassign PR: service error: ", err)
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
	h.logger.Info("reassign PR success: ", fiber.Map{
		"pull_request_id": req.PullRequestID,
		"replaced_by":     replacedBy,
	})

	return c.Status(fiber.StatusOK).JSON(response)
}
