package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/internal/logger"
)

// GetSchedules handles the request to get all schedules
func (h *Handler) GetSchedules(c *gin.Context) {
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetSchedules")
	log.Info("GetSchedules requested")

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

	var schedules interface{}

	// Apply filters
	if team != "" {
		log.WithField("team", team).Info("Getting schedules by team")
		schedules, err = h.schedulesRepo.FindByTeam(c.Request.Context(), team)
	} else if seasonStr != "" && weekStr != "" {
		log.WithFields(logrus.Fields{
			"season": season,
			"week":   week,
		}).Info("Getting schedules by week")
		schedules, err = h.schedulesRepo.FindByWeek(c.Request.Context(), season, week)
	} else {
		log.Info("Getting all schedules")
		schedules, err = h.schedulesRepo.FindAll(c.Request.Context())
	}

	if err != nil {
		log.WithError(err).Error("Failed to get schedules")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get schedules",
		})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// GetScheduleByID handles the request to get a schedule by ID
func (h *Handler) GetScheduleByID(c *gin.Context) {
	id := c.Param("id")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetScheduleByID").WithField("schedule_id", id)
	log.Info("GetScheduleByID requested")

	schedule, err := h.schedulesRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		log.WithError(err).Error("Failed to get schedule")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get schedule",
		})
		return
	}

	if schedule == nil {
		log.Info("Schedule not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Schedule not found",
		})
		return
	}

	log.Info("Schedule retrieved successfully")
	c.JSON(http.StatusOK, schedule)
}

// GetScheduleByGameKey handles the request to get a schedule by GameKey
func (h *Handler) GetScheduleByGameKey(c *gin.Context) {
	gameKey := c.Param("gameKey")
	log := logger.WithRequestContext(c.Request.Context()).WithField("component", "handlers.GetScheduleByGameKey").WithField("game_key", gameKey)
	log.Info("GetScheduleByGameKey requested")

	schedule, err := h.schedulesRepo.FindByGameKey(c.Request.Context(), gameKey)
	if err != nil {
		log.WithError(err).Error("Failed to get schedule")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get schedule",
		})
		return
	}

	if schedule == nil {
		log.Info("Schedule not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Schedule not found",
		})
		return
	}

	log.Info("Schedule retrieved successfully")
	c.JSON(http.StatusOK, schedule)
}
