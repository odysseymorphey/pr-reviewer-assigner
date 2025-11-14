package di

import (
	"log"
	"pr-reviwer-assigner/internal/config"
	"pr-reviwer-assigner/internal/domain/services"
	"pr-reviwer-assigner/internal/infrastructure/database"
	repo2 "pr-reviwer-assigner/internal/infrastructure/database/repository"

	"go.uber.org/zap"
)

type Container struct {
	prService   services.PRService
	teamService services.TeamService
	userService services.UserService

	logger *zap.SugaredLogger
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	zapLogger := zap.SugaredLogger{}

	//prLog := zapLogger.Named("prService")
	//teamLog := zapLogger.Named("teamService")
	//userLog := zapLogger.Named("userService")

	prrepo := repo2.NewPRRepository(db)
	teamrepo := repo2.NewTeamRepository(db)
	userrepo := repo2.NewUserRepository(db)

	prservice := services.NewPRService(prrepo)
	teamservice := services.NewTeamService(teamrepo)
	userservice := services.NewUserService(userrepo)

	return &Container{
		prService:   prservice,
		teamService: teamservice,
		userService: userservice,
		logger:      &zapLogger,
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

func (c *Container) GetNamedLogger(name string) *zap.SugaredLogger {
	return c.logger.Named(name)
}
