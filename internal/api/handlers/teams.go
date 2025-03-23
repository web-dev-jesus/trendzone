package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/web-dev-jesus/trendzone/internal/logger"
)

// GetTeams handles the request to get all teams
func (h *Handler) GetTeams(c *gin.Context) {
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetTeams")
	log.Info("GetTeams requested")

	teams, err := h.teamsRepo.FindAll(c.Request.Context())
	if err != nil {
		log.WithError(err).Error("Failed to get teams")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get teams",
		})
		return
	}

	log.WithField("count", len(teams)).Info("Teams retrieved successfully")
	c.JSON(http.StatusOK, teams)
}

// GetTeamByID handles the request to get a team by ID
func (h *Handler) GetTeamByID(c *gin.Context) {
	id := c.Param("id")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetTeamByID").WithField("team_id", id)
	log.Info("GetTeamByID requested")

	team, err := h.teamsRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		log.WithError(err).Error("Failed to get team")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get team",
		})
		return
	}

	if team == nil {
		log.Info("Team not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Team not found",
		})
		return
	}

	log.Info("Team retrieved successfully")
	c.JSON(http.StatusOK, team)
}

// GetTeamByKey handles the request to get a team by key (abbreviation)
func (h *Handler) GetTeamByKey(c *gin.Context) {
	key := c.Param("key")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetTeamByKey").WithField("team_key", key)
	log.Info("GetTeamByKey requested")

	team, err := h.teamsRepo.FindByKey(c.Request.Context(), key)
	if err != nil {
		log.WithError(err).Error("Failed to get team")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get team",
		})
		return
	}

	if team == nil {
		log.Info("Team not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Team not found",
		})
		return
	}

	log.Info("Team retrieved successfully")
	c.JSON(http.StatusOK, team)
}
