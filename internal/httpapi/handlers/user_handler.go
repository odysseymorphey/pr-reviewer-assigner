package handlers

import (
	"pr-reviwer-assigner/internal/domain/services"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SetIsActive(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (h *UserHandler) GetReview(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
