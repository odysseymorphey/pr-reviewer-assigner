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

type UserHandler struct {
	userService services.UserService
	logger      *zap.SugaredLogger
}

func NewUserHandler(userService services.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) SetIsActive(c fiber.Ctx) error {
	var req dto.SIARequest

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		h.logger.Error("failed to unmarshal body: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrInternal.Error(),
				Message: "internal server error",
			},
		})
	}

	resp, err := h.userService.SetIsActive(req)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrNotFound):
			h.logger.Error("failed to set IsActive: ", err.Error())

			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrNotFound.Error(),
					Message: "resource not found",
				},
			})
		default:
			h.logger.Error("failed to set IsActive: ", err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: dto.Error{
					Code:    errors2.ErrInternal.Error(),
					Message: "internal server error",
				},
			})
		}
	}

	h.logger.Info("SetIsActive success: ", resp)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *UserHandler) GetReview(c fiber.Ctx) error {
	userID := strings.TrimSpace(c.Params("user_id"))
	if userID == "" {
		h.logger.Error("empty user id")
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.Error{
				Code:    errors2.ErrBadRequest.Error(),
				Message: "user_id can't be empty",
			},
		})
	}

	// вот тут нужна проверка на отсутствуеюшего юзера
	prs, err := h.userService.GetReview(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{})
	}

	response := &dto.UserPR{
		ID:  userID,
		PRs: prs,
	}

	h.logger.Info("GetReview success: ", response)

	return c.Status(fiber.StatusOK).JSON(response)
}
