package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"

	"pr-reviwer-assigner/internal/config"
	"pr-reviwer-assigner/internal/di"
	"pr-reviwer-assigner/internal/httpapi"
)

type Server struct {
	app *fiber.App
	cfg *config.Config

	stopC chan os.Signal
}

func New(cfg *config.Config, c *di.Container) (*Server, error) {
	app := fiber.New()

	httpapi.RegisterRoutes(app, c)

	return &Server{
		app:   app,
		cfg:   cfg,
		stopC: make(chan os.Signal, 1),
	}, nil
}

func (s *Server) Run() {
	signal.Notify(s.stopC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.app.Listen(s.cfg.HTTPAddr); err != nil {
			log.Fatal(err)
		}
	}()

	<-s.stopC
	log.Println("Shutting down server...")

	s.stop()

	log.Println("Server stopped gracefully")
}

func (s *Server) stop() {
	s.app.Shutdown()
}
