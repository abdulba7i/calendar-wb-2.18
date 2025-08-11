package main

import (
	"wb-calendar/config"
	"wb-calendar/internal/calendar"
	"wb-calendar/internal/handler"
	"wb-calendar/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	logger.Init()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("No .env file found, using default values")
	}

	cfg := config.MustLoad()
	logger.Log.Infof("Server starting on %s", cfg.HTTPServer.Address)

	service := calendar.NewService()
	router := handler.InitRoute(service)

	if err := router.Run(cfg.HTTPServer.Address); err != nil {
		logger.Log.Fatalf("Failed to start server: %v", err)
	}
}
