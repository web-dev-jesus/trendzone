package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/api/handlers"
	"github.com/web-dev-jesus/trendzone/internal/api/routes"
	"github.com/web-dev-jesus/trendzone/internal/db/mongodb"
	"github.com/web-dev-jesus/trendzone/internal/db/mongodb/repositories"
	"github.com/web-dev-jesus/trendzone/internal/logger"
	"github.com/web-dev-jesus/trendzone/internal/sportsdata"
)

func main() {
	// Create a root context
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Setup logger
	logger.Setup(cfg.App.LogLevel)
	log := logrus.WithField("component", "main")

	log.Info("Starting NFL Stats Service")
	log.WithField("environment", cfg.App.Env).Info("Configuration loaded")

	// Connect to MongoDB
	mongoClient, err := mongodb.NewClient(ctx, &cfg.MongoDB)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to MongoDB")
	}

	// Create repositories
	teamsRepo := repositories.NewTeamsRepository(mongoClient.GetDatabase())
	playersRepo := repositories.NewPlayersRepository(mongoClient.GetDatabase())
	gamesRepo := repositories.NewGamesRepository(mongoClient.GetDatabase())
	standingsRepo := repositories.NewStandingsRepository(mongoClient.GetDatabase())
	schedulesRepo := repositories.NewSchedulesRepository(mongoClient.GetDatabase())

	// Create SportsData.io client and service
	sportsDataClient := sportsdata.NewClient(&cfg.SportsData)
	sportsDataService := sportsdata.NewService(
		sportsDataClient,
		teamsRepo,
		playersRepo,
		standingsRepo,
		schedulesRepo,
		gamesRepo,
	)

	// Create handler
	handler := handlers.NewHandler(
		cfg,
		teamsRepo,
		playersRepo,
		gamesRepo,
		standingsRepo,
		schedulesRepo,
		sportsDataService,
	)

	// Setup router
	router := routes.SetupRouter(cfg, handler)

	// Configure the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.WithField("port", cfg.App.Port).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("HTTP server failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	// Close MongoDB connection
	if err := mongoClient.Close(ctx); err != nil {
		log.WithError(err).Error("Failed to close MongoDB connection")
	}

	log.Info("Server exiting")
}
