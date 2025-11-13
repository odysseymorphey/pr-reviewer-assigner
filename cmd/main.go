package main

import (
	"log"
	"os"

	"pr-reviwer-assigner/internal/config"
	"pr-reviwer-assigner/internal/di"
	"pr-reviwer-assigner/internal/server"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config/config.json"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	container := di.NewContainer(cfg)

	app, err := server.New(cfg, container)
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
