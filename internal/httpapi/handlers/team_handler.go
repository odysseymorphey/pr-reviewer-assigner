package handlers

import (
	"pr-reviwer-assigner/internal/domain/services"

	"github.com/gofiber/fiber/v3"
)

type TeamHandler struct {
	teamService services.TeamService
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) Add(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (h *TeamHandler) Get(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
