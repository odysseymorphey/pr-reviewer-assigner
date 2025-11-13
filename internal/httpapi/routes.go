package httpapi

import (
	"pr-reviwer-assigner/internal/di"
	"pr-reviwer-assigner/internal/httpapi/handlers"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(r *fiber.App, c *di.Container) {
	teamHandler := handlers.NewTeamHandler(c.GetTeamService())
	userHandler := handlers.NewUserHandler(c.GetUserService())
	prHandler := handlers.NewPRHandler(c.GetPRService())

	// TEAM
	{
		r.Get("/team/get", teamHandler.Get)
		r.Post("/team/add", teamHandler.Add)
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
