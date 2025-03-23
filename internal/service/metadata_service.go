package service

import (
	"context"
	"time"

	"github.com/web-dev-jesus/trendzone/internal/api"
	"github.com/web-dev-jesus/trendzone/internal/db/repositories"
	"github.com/web-dev-jesus/trendzone/internal/models"
	"github.com/web-dev-jesus/trendzone/internal/util"
)

// MetadataService handles the synchronization and retrieval of NFL metadata
type MetadataService struct {
	apiClient      *api.Client
	teamRepo       *repositories.TeamRepository
	playerRepo     *repositories.PlayerRepository
	stadiumRepo    *repositories.StadiumRepository
	refereeRepo    *repositories.RefereeRepository
	byeWeekRepo    *repositories.ByeWeekRepository
	depthChartRepo *repositories.DepthChartRepository
	logger         *util.Logger
	lastSyncTime   map[string]time.Time
}

// NewMetadataService creates a new metadata service
func NewMetadataService(
	apiClient *api.Client,
	teamRepo *repositories.TeamRepository,
	playerRepo *repositories.PlayerRepository,
	stadiumRepo *repositories.StadiumRepository,
	refereeRepo *repositories.RefereeRepository,
	byeWeekRepo *repositories.ByeWeekRepository,
	depthChartRepo *repositories.DepthChartRepository,
	logger *util.Logger,
) *MetadataService {
	return &MetadataService{
		apiClient:    apiClient,
		teamRepo:     teamRepo,
		playerRepo:   playerRepo,
		stadiumRepo:  stadiumRepo,
		refereeRepo:  refereeRepo,
		byeWeekRepo:  byeWeekRepo,
		logger:       logger,
		lastSyncTime: make(map[string]time.Time),
	}
}

// SyncTeams fetches team data from the API and stores it in the database
func (s *MetadataService) SyncTeams(ctx context.Context) error {
	s.logger.Info("Starting team data synchronization")

	teams, err := s.apiClient.GetTeams(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch teams from API", "error", err)
		return err
	}

	for _, team := range teams {
		if err := s.teamRepo.UpsertTeam(ctx, team); err != nil {
			s.logger.Error("Failed to upsert team", "team", team.Key, "error", err)
			continue
		}
	}

	s.lastSyncTime["teams"] = time.Now()
	s.logger.Info("Team data synchronization completed", "count", len(teams))
	return nil
}

// SyncPlayers fetches player data from the API and stores it in the database
func (s *MetadataService) SyncPlayers(ctx context.Context) error {
	s.logger.Info("Starting player data synchronization")

	players, err := s.apiClient.GetPlayers(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch players from API", "error", err)
		return err
	}

	for _, player := range players {
		if err := s.playerRepo.UpsertPlayer(ctx, player); err != nil {
			s.logger.Error("Failed to upsert player", "playerID", player.PlayerID, "error", err)
			continue
		}
	}

	s.lastSyncTime["players"] = time.Now()
	s.logger.Info("Player data synchronization completed", "count", len(players))
	return nil
}

// SyncStadiums fetches stadium data from the API and stores it in the database
func (s *MetadataService) SyncStadiums(ctx context.Context) error {
	s.logger.Info("Starting stadium data synchronization")

	stadiums, err := s.apiClient.GetStadiums(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch stadiums from API", "error", err)
		return err
	}

	for _, stadium := range stadiums {
		if err := s.stadiumRepo.UpsertStadium(ctx, stadium); err != nil {
			s.logger.Error("Failed to upsert stadium", "stadiumID", stadium.StadiumID, "error", err)
			continue
		}
	}

	s.lastSyncTime["stadiums"] = time.Now()
	s.logger.Info("Stadium data synchronization completed", "count", len(stadiums))
	return nil
}

// SyncReferees fetches referee data from the API and stores it in the database
func (s *MetadataService) SyncReferees(ctx context.Context) error {
	s.logger.Info("Starting referee data synchronization")

	referees, err := s.apiClient.GetReferees(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch referees from API", "error", err)
		return err
	}

	for _, referee := range referees {
		if err := s.refereeRepo.UpsertReferee(ctx, referee); err != nil {
			s.logger.Error("Failed to upsert referee", "refereeID", referee.RefereeID, "error", err)
			continue
		}
	}

	s.lastSyncTime["referees"] = time.Now()
	s.logger.Info("Referee data synchronization completed", "count", len(referees))
	return nil
}

// SyncByeWeeks fetches bye week data for a specific season and stores it in the database
func (s *MetadataService) SyncByeWeeks(ctx context.Context, season string) error {
	s.logger.Info("Starting bye week data synchronization", "season", season)

	byeWeeks, err := s.apiClient.GetByeWeeks(ctx, season)
	if err != nil {
		s.logger.Error("Failed to fetch bye weeks from API", "season", season, "error", err)
		return err
	}

	for _, byeWeek := range byeWeeks {
		if err := s.byeWeekRepo.UpsertByeWeek(ctx, byeWeek); err != nil {
			s.logger.Error("Failed to upsert bye week", "team", byeWeek.Team, "season", season, "error", err)
			continue
		}
	}

	s.lastSyncTime["byeWeeks_"+season] = time.Now()
	s.logger.Info("Bye week data synchronization completed", "season", season, "count", len(byeWeeks))
	return nil
}

// SyncDepthCharts fetches team depth chart data and stores it in the database
func (s *MetadataService) SyncDepthCharts(ctx context.Context) error {
	s.logger.Info("Starting depth chart data synchronization")

	depthCharts, err := s.apiClient.GetDepthCharts(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch depth charts from API", "error", err)
		return err
	}

	for _, teamDepthChart := range depthCharts {
		depthChart := models.DepthChart{
			TeamID:       teamDepthChart.TeamID,
			Team:         "", // You might need to fetch the team abbreviation from the team repository
			Offense:      make([]models.DepthChartPosition, 0),
			Defense:      make([]models.DepthChartPosition, 0),
			SpecialTeams: make([]models.DepthChartPosition, 0),
			LastUpdated:  time.Now(),
		}

		// Get the team abbreviation if needed
		team, err := s.teamRepo.GetTeamByID(ctx, teamDepthChart.TeamID)
		if err == nil && team != nil {
			depthChart.Team = team.Key
		} else {
			s.logger.Warn("Could not find team for depth chart", "teamID", teamDepthChart.TeamID)
		}

		// Process offensive positions
		for _, position := range teamDepthChart.Offense {
			depthChart.Offense = append(depthChart.Offense, models.DepthChartPosition{
				DepthChartID:     position.DepthChartID,
				TeamID:           position.TeamID,
				PlayerID:         position.PlayerID,
				Name:             position.Name,
				PositionCategory: position.PositionCategory,
				Position:         position.Position,
				DepthOrder:       position.DepthOrder,
			})
		}

		// Process defensive positions
		for _, position := range teamDepthChart.Defense {
			depthChart.Defense = append(depthChart.Defense, models.DepthChartPosition{
				DepthChartID:     position.DepthChartID,
				TeamID:           position.TeamID,
				PlayerID:         position.PlayerID,
				Name:             position.Name,
				PositionCategory: position.PositionCategory,
				Position:         position.Position,
				DepthOrder:       position.DepthOrder,
			})
		}

		// Process special teams positions
		for _, position := range teamDepthChart.SpecialTeams {
			depthChart.SpecialTeams = append(depthChart.SpecialTeams, models.DepthChartPosition{
				DepthChartID:     position.DepthChartID,
				TeamID:           position.TeamID,
				PlayerID:         position.PlayerID,
				Name:             position.Name,
				PositionCategory: position.PositionCategory,
				Position:         position.Position,
				DepthOrder:       position.DepthOrder,
			})
		}

		// Upsert the depth chart to the database
		if err := s.depthChartRepo.UpsertDepthChart(ctx, &depthChart); err != nil {
			s.logger.Error("Failed to upsert depth chart", "teamID", depthChart.TeamID, "error", err)
			continue
		}
	}

	s.lastSyncTime["depthCharts"] = time.Now()
	s.logger.Info("Depth chart data synchronization completed", "count", len(depthCharts))
	return nil
}

// SyncAllMetadata synchronizes all available metadata
func (s *MetadataService) SyncAllMetadata(ctx context.Context, season string) error {
	s.logger.Info("Starting full metadata synchronization")

	if err := s.SyncTeams(ctx); err != nil {
		s.logger.Error("Teams sync failed during full sync", "error", err)
		// Continue with other syncs even if one fails
	}

	if err := s.SyncPlayers(ctx); err != nil {
		s.logger.Error("Players sync failed during full sync", "error", err)
	}

	if err := s.SyncStadiums(ctx); err != nil {
		s.logger.Error("Stadiums sync failed during full sync", "error", err)
	}

	if err := s.SyncReferees(ctx); err != nil {
		s.logger.Error("Referees sync failed during full sync", "error", err)
	}

	if err := s.SyncByeWeeks(ctx, season); err != nil {
		s.logger.Error("Bye weeks sync failed during full sync", "season", season, "error", err)
	}

	if err := s.SyncDepthCharts(ctx); err != nil {
		s.logger.Error("Depth charts sync failed during full sync", "error", err)
	}

	s.logger.Info("Full metadata synchronization completed")
	return nil
}

// GetLastSyncTime returns the last time a specific entity type was synchronized
func (s *MetadataService) GetLastSyncTime(entityType string) time.Time {
	return s.lastSyncTime[entityType]
}

// GetTeamByID retrieves a team by its ID from the database
func (s *MetadataService) GetTeamByID(ctx context.Context, teamID int) (*models.Team, error) {
	return s.teamRepo.GetTeamByID(ctx, teamID)
}

// GetTeamByKey retrieves a team by its key (abbreviation) from the database
func (s *MetadataService) GetTeamByKey(ctx context.Context, key string) (*models.Team, error) {
	return s.teamRepo.GetTeamByKey(ctx, key)
}

// GetPlayerByID retrieves a player by ID from the database
func (s *MetadataService) GetPlayerByID(ctx context.Context, playerID int) (*models.Player, error) {
	return s.playerRepo.GetPlayerByID(ctx, playerID)
}

// GetPlayersByTeam retrieves all players for a specific team from the database
func (s *MetadataService) GetPlayersByTeam(ctx context.Context, teamKey string) ([]*models.Player, error) {
	return s.playerRepo.GetPlayersByTeam(ctx, teamKey)
}

// GetStadiumByID retrieves a stadium by ID from the database
func (s *MetadataService) GetStadiumByID(ctx context.Context, stadiumID int) (*models.Stadium, error) {
	return s.stadiumRepo.GetStadiumByID(ctx, stadiumID)
}

// GetAllReferees retrieves all referees from the database
func (s *MetadataService) GetAllReferees(ctx context.Context) ([]*models.Referee, error) {
	return s.refereeRepo.GetAllReferees(ctx)
}

// GetByeWeeksForSeason retrieves all bye weeks for a specific season
func (s *MetadataService) GetByeWeeksForSeason(ctx context.Context, season string) ([]*models.ByeWeek, error) {
	return s.byeWeekRepo.GetByeWeeksForSeason(ctx, season)
}

// GetTeamByeWeek retrieves the bye week for a specific team and season
func (s *MetadataService) GetTeamByeWeek(ctx context.Context, teamKey string, season string) (*models.ByeWeek, error) {
	return s.byeWeekRepo.GetTeamByeWeek(ctx, teamKey, season)
}

// GetTeamDepthChart retrieves a team's depth chart from the database
func (s *MetadataService) GetTeamDepthChart(ctx context.Context, teamID int) (*models.DepthChart, error) {
	return s.depthChartRepo.GetDepthChartByTeamID(ctx, teamID)
}

// GetAllDepthCharts retrieves all team depth charts from the database
func (s *MetadataService) GetAllDepthCharts(ctx context.Context) ([]*models.DepthChart, error) {
	return s.depthChartRepo.GetAllDepthCharts(ctx)
}
