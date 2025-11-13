package handlers

import (
	"encoding/json"
	"pr-reviwer-assigner/internal/domain/dto"
	"pr-reviwer-assigner/internal/domain/services"

	"github.com/gofiber/fiber/v3"
)

const (
	ErrNotFound = "NOT_FOUND"
	ErrExists   = "PR_EXISTS"
)

type PRHandler struct {
	service services.PRService
}

func NewPRHandler(service services.PRService) *PRHandler {
	return &PRHandler{
		service: service,
	}
}

func (h *PRHandler) CreatePR(c fiber.Ctx) error {
	var prReq dto.PRRequest

	err := json.Unmarshal(c.Body(), &prReq)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
	}

	h.service.Create(prReq)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (h *PRHandler) MergePR(c fiber.Ctx) error {
	return nil
}

func (h *PRHandler) ReassignViewer(c fiber.Ctx) error {
	return nil
}
