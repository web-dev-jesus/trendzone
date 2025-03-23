package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/db/mongodb/repositories"
	"github.com/web-dev-jesus/trendzone/internal/logger"
	"github.com/web-dev-jesus/trendzone/internal/sportsdata"
)

type Handler struct {
	config            *config.Config
	teamsRepo         *repositories.TeamsRepository
	playersRepo       *repositories.PlayersRepository
	gamesRepo         *repositories.GamesRepository
	standingsRepo     *repositories.StandingsRepository
	schedulesRepo     *repositories.SchedulesRepository
	sportsDataService *sportsdata.Service
}

func NewHandler(
	config *config.Config,
	teamsRepo *repositories.TeamsRepository,
	playersRepo *repositories.PlayersRepository,
	gamesRepo *repositories.GamesRepository,
	standingsRepo *repositories.StandingsRepository,
	schedulesRepo *repositories.SchedulesRepository,
	sportsDataService *sportsdata.Service,
) *Handler {
	return &Handler{
		config:            config,
		teamsRepo:         teamsRepo,
		playersRepo:       playersRepo,
		gamesRepo:         gamesRepo,
		standingsRepo:     standingsRepo,
		schedulesRepo:     schedulesRepo,
		sportsDataService: sportsDataService,
	}
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.HealthCheck")
	log.Info("Health check requested")

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
	})
}

// SyncData handles the data synchronization endpoint
func (h *Handler) SyncData(c *gin.Context) {
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.SyncData")

	// Extract season from query parameters, default to current year
	season := c.DefaultQuery("season", "2023")

	log.WithField("season", season).Info("Data sync requested")

	// Start the sync process asynchronously
	go func() {
		ctx := c.Request.Context()
		if err := h.sportsDataService.SyncAll(ctx, season); err != nil {
			log.WithError(err).Error("Failed to sync data")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Data synchronization started",
		"season":  season,
	})
}
