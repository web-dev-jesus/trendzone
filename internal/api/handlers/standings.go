package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/internal/logger"
)

// GetStandings handles the request to get all standings
func (h *Handler) GetStandings(c *gin.Context) {
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetStandings")
	log.Info("GetStandings requested")

	// Check for filters
	conference := c.Query("conference")
	division := c.Query("division")

	var standings interface{}
	var err error

	// Apply filters
	if conference != "" && division != "" {
		log.WithFields(logrus.Fields{
			"conference": conference,
			"division":   division,
		}).Info("Getting standings by division")
		standings, err = h.standingsRepo.FindByDivision(c.Request.Context(), conference, division)
	} else {
		log.Info("Getting all standings")
		standings, err = h.standingsRepo.FindAll(c.Request.Context())
	}

	if err != nil {
		log.WithError(err).Error("Failed to get standings")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get standings",
		})
		return
	}

	c.JSON(http.StatusOK, standings)
}

// GetStandingByID handles the request to get a standing by ID
func (h *Handler) GetStandingByID(c *gin.Context) {
	id := c.Param("id")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetStandingByID").WithField("standing_id", id)
	log.Info("GetStandingByID requested")

	standing, err := h.standingsRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		log.WithError(err).Error("Failed to get standing")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get standing",
		})
		return
	}

	if standing == nil {
		log.Info("Standing not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Standing not found",
		})
		return
	}

	log.Info("Standing retrieved successfully")
	c.JSON(http.StatusOK, standing)
}

// GetStandingByTeam handles the request to get a standing by team
func (h *Handler) GetStandingByTeam(c *gin.Context) {
	team := c.Param("team")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetStandingByTeam").WithField("team", team)
	log.Info("GetStandingByTeam requested")

	standing, err := h.standingsRepo.FindByTeam(c.Request.Context(), team)
	if err != nil {
		log.WithError(err).Error("Failed to get standing")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get standing",
		})
		return
	}

	if standing == nil {
		log.Info("Standing not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Standing not found",
		})
		return
	}

	log.Info("Standing retrieved successfully")
	c.JSON(http.StatusOK, standing)
}
