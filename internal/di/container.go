package di

import (
	"log"
	"pr-reviwer-assigner/internal/config"
	"pr-reviwer-assigner/internal/domain/services"
	"pr-reviwer-assigner/internal/infrastructure/database"
	repo2 "pr-reviwer-assigner/internal/infrastructure/database/repository"
)

type Container struct {
	prService   services.PRService
	teamService services.TeamService
	userService services.UserService
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	prrepo := repo2.NewPRRepository(db)
	prservice := services.NewPRService(prrepo)

	teamrepo := repo2.NewTeamRepository(db)
	teamservice := services.NewTeamService(teamrepo)

	userrepo := repo2.NewUserRepository(db)
	userservice := services.NewUserService(userrepo)

	return &Container{
		prService:   prservice,
		teamService: teamservice,
		userService: userservice,
	}
}

func (c *Container) GetPRService() services.PRService {
	return c.prService
}

func (c *Container) GetTeamService() services.TeamService {
	return c.teamService
}

func (c *Container) GetUserService() services.UserService {
	return c.userService
}
