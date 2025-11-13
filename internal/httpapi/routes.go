package httpapi

import (
	"pr-reviwer-assigner/internal/di"
	"pr-reviwer-assigner/internal/httpapi/handlers"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(r *fiber.App, c *di.Container) {
	prHandler := handlers.NewPRHandler(c.GetPRService())

	// TEAM
	{
		r.Get("/team/get", func(ctx fiber.Ctx) error {
			return nil
		})
		r.Post("/team/add", func(ctx fiber.Ctx) error {
			return nil
		})
	}

	// USERS
	{
		r.Post("/users/setIsActive", func(ctx fiber.Ctx) error {
			return nil
		})
		r.Get("/users/getReview", func(ctx fiber.Ctx) error {
			return nil
		})
	}

	// PR
	{
		r.Post("/pullRequest/create", prHandler.CreatePR)
		r.Post("/pullRequest/merge", prHandler.MergePR)
		r.Post("/pullRequest/reassign", prHandler.ReassignViewer)
	}
}
