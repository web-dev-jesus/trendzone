package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/api/handlers"
	"github.com/web-dev-jesus/trendzone/internal/api/middleware"
)

func SetupRouter(cfg *config.Config, handler *handlers.Handler) *gin.Engine {
	r := gin.New()

	// Apply middleware
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// API documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", handler.HealthCheck)

	// API routes
	apiV1 := r.Group("/api/v1")
	{
		// Public routes
		// Teams
		apiV1.GET("/teams", handler.GetTeams)
		apiV1.GET("/teams/:id", handler.GetTeamByID)
		apiV1.GET("/teams/key/:key", handler.GetTeamByKey)

		// Players
		apiV1.GET("/players", handler.GetPlayers)
		apiV1.GET("/players/:id", handler.GetPlayerByID)
		apiV1.GET("/players/pid/:playerID", handler.GetPlayerByPlayerID)

		// Games
		apiV1.GET("/games", handler.GetGames)
		apiV1.GET("/games/:id", handler.GetGameByID)
		apiV1.GET("/games/key/:gameKey", handler.GetGameByGameKey)

		// Standings
		apiV1.GET("/standings", handler.GetStandings)
		apiV1.GET("/standings/:id", handler.GetStandingByID)
		apiV1.GET("/standings/team/:team", handler.GetStandingByTeam)

		// Schedules
		apiV1.GET("/schedules", handler.GetSchedules)
		apiV1.GET("/schedules/:id", handler.GetScheduleByID)
		apiV1.GET("/schedules/key/:gameKey", handler.GetScheduleByGameKey)

		// Protected routes (require authentication)
		adminRoutes := apiV1.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(&cfg.App))
		{
			// Data sync
			adminRoutes.POST("/sync", handler.SyncData)
		}
	}

	return r
}
