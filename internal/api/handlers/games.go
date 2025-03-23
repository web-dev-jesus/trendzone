package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/internal/logger"
)

// GetGames handles the request to get all games
func (h *Handler) GetGames(c *gin.Context) {
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetGames")
	log.Info("GetGames requested")

	// Check for filters
	team := c.Query("team")
	seasonStr := c.Query("season")
	weekStr := c.Query("week")

	var season, week int
	var err error

	if seasonStr != "" {
		season, err = strconv.Atoi(seasonStr)
		if err != nil {
			log.WithError(err).Error("Invalid season format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid season format",
			})
			return
		}
	}

	if weekStr != "" {
		week, err = strconv.Atoi(weekStr)
		if err != nil {
			log.WithError(err).Error("Invalid week format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid week format",
			})
			return
		}
	}

	var games interface{}

	// Apply filters
	if team != "" {
		log.WithField("team", team).Info("Getting games by team")
		games, err = h.gamesRepo.FindByTeam(c.Request.Context(), team)
	} else if seasonStr != "" && weekStr != "" {
		log.WithFields(logrus.Fields{
			"season": season,
			"week":   week,
		}).Info("Getting games by week")
		games, err = h.gamesRepo.FindByWeek(c.Request.Context(), season, week)
	} else {
		log.Info("Getting all games")
		games, err = h.gamesRepo.FindAll(c.Request.Context())
	}

	if err != nil {
		log.WithError(err).Error("Failed to get games")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get games",
		})
		return
	}

	c.JSON(http.StatusOK, games)
}

// GetGameByID handles the request to get a game by ID
func (h *Handler) GetGameByID(c *gin.Context) {
	id := c.Param("id")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetGameByID").WithField("game_id", id)
	log.Info("GetGameByID requested")

	game, err := h.gamesRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		log.WithError(err).Error("Failed to get game")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get game",
		})
		return
	}

	if game == nil {
		log.Info("Game not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	log.Info("Game retrieved successfully")
	c.JSON(http.StatusOK, game)
}

// GetGameByGameKey handles the request to get a game by GameKey
func (h *Handler) GetGameByGameKey(c *gin.Context) {
	gameKey := c.Param("gameKey")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetGameByGameKey").WithField("game_key", gameKey)
	log.Info("GetGameByGameKey requested")

	game, err := h.gamesRepo.FindByGameKey(c.Request.Context(), gameKey)
	if err != nil {
		log.WithError(err).Error("Failed to get game")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get game",
		})
		return
	}

	if game == nil {
		log.Info("Game not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
		})
		return
	}

	log.Info("Game retrieved successfully")
	c.JSON(http.StatusOK, game)
}
