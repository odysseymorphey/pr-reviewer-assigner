package httpapi

import (
	"pr-reviwer-assigner/internal/di"
	"pr-reviwer-assigner/internal/httpapi/docs"
	"pr-reviwer-assigner/internal/httpapi/handlers"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(r *fiber.App, c *di.Container) {
	teamHandler := handlers.NewTeamHandler(c.GetTeamService(), c.GetNamedLogger("teamHandler"))
	userHandler := handlers.NewUserHandler(c.GetUserService(), c.GetNamedLogger("userHandler"))
	prHandler := handlers.NewPRHandler(c.GetPRService(), c.GetNamedLogger("prHandler"))
	docs.RegisterRoutes(r)

	// HEALTH
	{
		r.Get("/health", func(c fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"health": "ok",
			})
		})
	}

	// TEAM
	{
		r.Get("/team/get", teamHandler.Get)
		r.Post("/team/add", teamHandler.Add)
		r.Post("/team/deactivateMembers", teamHandler.DeactivateMembers)
	}

	// USERS
	{
		r.Post("/users/setIsActive", userHandler.SetIsActive)
		r.Get("/users/getReview", userHandler.GetReview)
	}

	// PR
	{
		r.Post("/pullRequest/create", prHandler.CreatePR)
		r.Post("/pullRequest/merge", prHandler.MergePR)
		r.Post("/pullRequest/reassign", prHandler.ReassignViewer)
	}
}
