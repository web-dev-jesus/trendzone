package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/web-dev-jesus/trendzone/internal/db/models"
	"github.com/web-dev-jesus/trendzone/internal/logger"
)

// GetPlayers handles the request to get all players
func (h *Handler) GetPlayers(c *gin.Context) {
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetPlayers")
	log.Info("GetPlayers requested")

	// Check if team filter is provided
	team := c.Query("team")

	var players []models.Player
	var err error

	if team != "" {
		log.WithField("team", team).Info("Getting players by team")
		players, err = h.playersRepo.FindByTeam(c.Request.Context(), team)
	} else {
		log.Info("Getting all players")
		players, err = h.playersRepo.FindAll(c.Request.Context())
	}

	if err != nil {
		log.WithError(err).Error("Failed to get players")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get players",
		})
		return
	}

	log.WithField("count", len(players)).Info("Players retrieved successfully")
	c.JSON(http.StatusOK, players)
}

// GetPlayerByID handles the request to get a player by ID
func (h *Handler) GetPlayerByID(c *gin.Context) {
	id := c.Param("id")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetPlayerByID").WithField("player_id", id)
	log.Info("GetPlayerByID requested")

	player, err := h.playersRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		log.WithError(err).Error("Failed to get player")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get player",
		})
		return
	}

	if player == nil {
		log.Info("Player not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Player not found",
		})
		return
	}

	log.Info("Player retrieved successfully")
	c.JSON(http.StatusOK, player)
}

// GetPlayerByPlayerID handles the request to get a player by PlayerID (from SportsData.io)
func (h *Handler) GetPlayerByPlayerID(c *gin.Context) {
	playerIDStr := c.Param("playerID")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetPlayerByPlayerID").WithField("player_id", playerIDStr)
	log.Info("GetPlayerByPlayerID requested")

	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid player ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid player ID format",
		})
		return
	}

	player, err := h.playersRepo.FindByPlayerID(c.Request.Context(), playerID)
	if err != nil {
		log.WithError(err).Error("Failed to get player")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get player",
		})
		return
	}

	if player == nil {
		log.Info("Player not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Player not found",
		})
		return
	}

	log.Info("Player retrieved successfully")
	c.JSON(http.StatusOK, player)
}
