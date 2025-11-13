package di

import (
	"log"
	"pr-reviwer-assigner/internal/config"
	"pr-reviwer-assigner/internal/domain/services"
	"pr-reviwer-assigner/internal/infrastructure/database"
	repo2 "pr-reviwer-assigner/internal/infrastructure/database/repository"
)

type Container struct {
	prService services.PRService
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	prrepo := repo2.NewPRRepository(db)
	prservice := services.NewPRService(prrepo)

	return &Container{
		prService: prservice,
	}
}

func (c *Container) GetPRService() services.PRService {
	return c.prService
}
