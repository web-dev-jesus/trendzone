package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/api"
	"github.com/web-dev-jesus/trendzone/internal/db"
	"github.com/web-dev-jesus/trendzone/internal/db/repositories"
	"github.com/web-dev-jesus/trendzone/internal/service"
	"github.com/web-dev-jesus/trendzone/internal/util"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := util.NewLogger(cfg.LogLevel)
	logger.Info("Starting NFL Data Sync Application")

	// Connect to MongoDB
	database, err := db.Connect(cfg, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
		os.Exit(1)
	}
	defer db.Disconnect(logger)

	// Initialize API client
	apiClient := api.NewClient(cfg, logger)
	defer apiClient.Close()

	// Initialize repositories
	teamsRepo := repositories.NewTeamsRepository(database, logger)
	playersRepo := repositories.NewPlayersRepository(database, logger)
	gamesRepo := repositories.NewGamesRepository(database, logger)
	standingsRepo := repositories.NewStandingsRepository(database, logger)
	schedulesRepo := repositories.NewSchedulesRepository(database, logger)
	playerGameStatsRepo := repositories.NewPlayerGameStatsRepository(database, logger)
	playByPlayRepo := repositories.NewPlayByPlayRepository(database, logger)
	stadiumsRepo := repositories.NewStadiumsRepository(database, logger)
	depthChartsRepo := repositories.NewDepthChartsRepository(database, logger)
	refereesRepo := repositories.NewRefereesRepository(database, logger)
	byeWeeksRepo := repositories.NewByeWeeksRepository(database, logger)
	metadataRepo := repositories.NewMetadataRepository(database, logger)

	// Initialize sync service
	syncService := service.NewSyncService(
		cfg,
		apiClient,
		teamsRepo,
		playersRepo,
		gamesRepo,
		standingsRepo,
		schedulesRepo,
		playerGameStatsRepo,
		playByPlayRepo,
		stadiumsRepo,
		depthChartsRepo,
		refereesRepo,
		byeWeeksRepo,
		metadataRepo,
		logger,
	)

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		logger.Info("Received shutdown signal")
		cancel()
	}()

	// Run sync process
	if err := syncService.SyncAllData(); err != nil {
		logger.Error(fmt.Sprintf("Error syncing data: %v", err))
		os.Exit(1)
	}

	// If you want to run as a scheduled task
	ticker := time.NewTicker(24 * time.Hour) // Daily update
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Info("Running scheduled data sync")
			if err := syncService.SyncAllData(); err != nil {
				logger.Error(fmt.Sprintf("Error in scheduled sync: %v", err))
			}
		case <-ctx.Done():
			logger.Info("Shutting down application")
			return
		}
	}
}
