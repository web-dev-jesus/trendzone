package service

import (
	"fmt"
	"time"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/api"
	"github.com/web-dev-jesus/trendzone/internal/db/repositories"
	"github.com/web-dev-jesus/trendzone/internal/models"
	"github.com/web-dev-jesus/trendzone/internal/util"
)

// SyncService orchestrates data syncing from SportsData.io to MongoDB
type SyncService struct {
	config              *config.Config
	apiClient           *api.Client
	teamsRepo           *repositories.TeamsRepository
	playersRepo         *repositories.PlayersRepository
	gamesRepo           *repositories.GamesRepository
	standingsRepo       *repositories.StandingsRepository
	schedulesRepo       *repositories.SchedulesRepository
	playerGameStatsRepo *repositories.PlayerGameStatsRepository
	playByPlayRepo      *repositories.PlayByPlayRepository
	stadiumsRepo        *repositories.StadiumsRepository
	depthChartsRepo     *repositories.DepthChartsRepository
	refereesRepo        *repositories.RefereesRepository
	byeWeeksRepo        *repositories.ByeWeeksRepository
	metadataRepo        *repositories.MetadataRepository
	logger              *util.Logger
}

// NewSyncService creates a new sync service
func NewSyncService(
	config *config.Config,
	apiClient *api.Client,
	teamsRepo *repositories.TeamsRepository,
	playersRepo *repositories.PlayersRepository,
	gamesRepo *repositories.GamesRepository,
	standingsRepo *repositories.StandingsRepository,
	schedulesRepo *repositories.SchedulesRepository,
	playerGameStatsRepo *repositories.PlayerGameStatsRepository,
	playByPlayRepo *repositories.PlayByPlayRepository,
	stadiumsRepo *repositories.StadiumsRepository,
	depthChartsRepo *repositories.DepthChartsRepository,
	refereesRepo *repositories.RefereesRepository,
	byeWeeksRepo *repositories.ByeWeeksRepository,
	metadataRepo *repositories.MetadataRepository,
	logger *util.Logger,
) *SyncService {
	return &SyncService{
		config:              config,
		apiClient:           apiClient,
		teamsRepo:           teamsRepo,
		playersRepo:         playersRepo,
		gamesRepo:           gamesRepo,
		standingsRepo:       standingsRepo,
		schedulesRepo:       schedulesRepo,
		playerGameStatsRepo: playerGameStatsRepo,
		playByPlayRepo:      playByPlayRepo,
		stadiumsRepo:        stadiumsRepo,
		depthChartsRepo:     depthChartsRepo,
		refereesRepo:        refereesRepo,
		byeWeeksRepo:        byeWeeksRepo,
		metadataRepo:        metadataRepo,
		logger:              logger,
	}
}

// SyncAllData syncs all NFL data from SportsData.io
func (s *SyncService) SyncAllData() error {
	s.logger.Info("Starting NFL data sync process")

	// Sync reference data that doesn't change often
	if err := s.SyncTeams(); err != nil {
		s.logger.Error("Error syncing teams: " + err.Error())
	}

	if err := s.SyncStadiums(); err != nil {
		s.logger.Error("Error syncing stadiums: " + err.Error())
	}

	if err := s.SyncReferees(); err != nil {
		s.logger.Error("Error syncing referees: " + err.Error())
	}

	if err := s.SyncByeWeeks(); err != nil {
		s.logger.Error("Error syncing bye weeks: " + err.Error())
	}

	// Sync core data
	if err := s.SyncPlayers(); err != nil {
		s.logger.Error("Error syncing players: " + err.Error())
	}

	if err := s.SyncDepthCharts(); err != nil {
		s.logger.Error("Error syncing depth charts: " + err.Error())
	}

	if err := s.SyncStandings(); err != nil {
		s.logger.Error("Error syncing standings: " + err.Error())
	}

	if err := s.SyncSchedules(); err != nil {
		s.logger.Error("Error syncing schedules: " + err.Error())
	}

	// Sync game data for each week
	currentWeek, err := s.getCurrentWeek()
	if err != nil {
		s.logger.Error("Error getting current week: " + err.Error())
		currentWeek = 17 // Default to maximum regular season week
	}

	for week := 1; week <= currentWeek; week++ {
		if err := s.SyncGames(week); err != nil {
			s.logger.Error(fmt.Sprintf("Error syncing games for week %d: %s", week, err.Error()))
		}

		if err := s.SyncPlayerGameStats(week); err != nil {
			s.logger.Error(fmt.Sprintf("Error syncing player game stats for week %d: %s", week, err.Error()))
		}

		if err := s.SyncPlayByPlay(week); err != nil {
			s.logger.Error(fmt.Sprintf("Error syncing play by play for week %d: %s", week, err.Error()))
		}
	}

	s.logger.Info("Completed NFL data sync process")
	return nil
}

// SyncTeams syncs teams data
func (s *SyncService) SyncTeams() error {
	// Check if update is needed
	needsUpdate, err := s.metadataRepo.IsUpdateNeeded("last_api_call_TeamsBasic", 168) // Weekly update
	if err != nil {
		return err
	}

	if !needsUpdate {
		s.logger.Info("Teams data is up to date")
		return nil
	}

	s.logger.Info("Syncing teams data")

	// Fetch teams data
	var teams []models.Team
	if err := s.apiClient.FetchData(api.TeamsEndpoint, &teams); err != nil {
		s.metadataRepo.RecordAPICall("last_api_call_TeamsBasic", api.TeamsEndpoint, "error", err.Error())
		return err
	}

	// Upsert teams data
	count, err := s.teamsRepo.UpsertTeams(teams)
	if err != nil {
		s.metadataRepo.RecordAPICall("last_api_call_TeamsBasic", api.TeamsEndpoint, "error", err.Error())
		return err
	}

	// Record metadata
	s.metadataRepo.RecordAPICall("last_api_call_TeamsBasic", api.TeamsEndpoint, "success", fmt.Sprintf("Updated %d teams", count))
	s.logger.Info(fmt.Sprintf("Successfully updated %d teams", count))

	return nil
}

// Additional sync methods would follow the same pattern

// getCurrentWeek determines the current NFL week
func (s *SyncService) getCurrentWeek() (int, error) {
	// Try to get from schedules
	currentDate := time.Now()
	upcomingGame, err := s.schedulesRepo.GetNextGame(currentDate)
	if err != nil {
		return 0, err
	}

	if upcomingGame != nil {
		return upcomingGame.Week, nil
	}

	// Fallback to a conservative estimate
	return 17, nil // Max regular season week
}
