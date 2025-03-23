package sportsdata

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/internal/db/mongodb/repositories"
	"github.com/web-dev-jesus/trendzone/internal/logger"
)

type Service struct {
	client        *Client
	teamsRepo     *repositories.TeamsRepository
	playersRepo   *repositories.PlayersRepository
	standingsRepo *repositories.StandingsRepository
	schedulesRepo *repositories.SchedulesRepository
	gamesRepo     *repositories.GamesRepository
}

func NewService(
	client *Client,
	teamsRepo *repositories.TeamsRepository,
	playersRepo *repositories.PlayersRepository,
	standingsRepo *repositories.StandingsRepository,
	schedulesRepo *repositories.SchedulesRepository,
	gamesRepo *repositories.GamesRepository,
) *Service {
	return &Service{
		client:        client,
		teamsRepo:     teamsRepo,
		playersRepo:   playersRepo,
		standingsRepo: standingsRepo,
		schedulesRepo: schedulesRepo,
		gamesRepo:     gamesRepo,
	}
}

// SyncTeams fetches teams from SportsData.io API and stores them in the database
func (s *Service) SyncTeams(ctx context.Context) error {
	log := logger.WithRequestContext(ctx).WithField("component", "sportsdata_service.SyncTeams")
	log.Info("Syncing teams from SportsData.io API to database")

	teams, err := s.client.GetTeams(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to fetch teams from API")
		return err
	}

	log.WithField("count", len(teams)).Info("Upserting teams in database")

	successCount := 0
	for _, team := range teams {
		_, err := s.teamsRepo.UpsertByTeamID(ctx, &team)
		if err != nil {
			log.WithFields(logrus.Fields{
				"team_id":  team.TeamID,
				"team_key": team.Key,
				"error":    err.Error(),
			}).Error("Failed to upsert team")
			continue
		}
		successCount++
	}

	log.WithFields(logrus.Fields{
		"success_count": successCount,
		"total_count":   len(teams),
	}).Info("Teams sync completed")

	return nil
}

// SyncPlayers fetches players from SportsData.io API and stores them in the database
func (s *Service) SyncPlayers(ctx context.Context) error {
	log := logger.WithRequestContext(ctx).WithField("component", "sportsdata_service.SyncPlayers")
	log.Info("Syncing players from SportsData.io API to database")

	players, err := s.client.GetPlayers(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to fetch players from API")
		return err
	}

	log.WithField("count", len(players)).Info("Upserting players in database")

	successCount := 0
	for _, player := range players {
		_, err := s.playersRepo.UpsertByPlayerID(ctx, &player)
		if err != nil {
			log.WithFields(logrus.Fields{
				"player_id":   player.PlayerID,
				"player_name": player.Name,
				"error":       err.Error(),
			}).Error("Failed to upsert player")
			continue
		}
		successCount++
	}

	log.WithFields(logrus.Fields{
		"success_count": successCount,
		"total_count":   len(players),
	}).Info("Players sync completed")

	return nil
}

// SyncStandings fetches standings from SportsData.io API and stores them in the database
func (s *Service) SyncStandings(ctx context.Context, season string) error {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "sportsdata_service.SyncStandings",
		"season":    season,
	})
	log.Info("Syncing standings from SportsData.io API to database")

	standings, err := s.client.GetStandings(ctx, season)
	if err != nil {
		log.WithError(err).Error("Failed to fetch standings from API")
		return err
	}

	log.WithField("count", len(standings)).Info("Upserting standings in database")

	successCount := 0
	for _, standing := range standings {
		_, err := s.standingsRepo.UpsertByTeamAndSeason(ctx, &standing)
		if err != nil {
			log.WithFields(logrus.Fields{
				"team":  standing.Team,
				"name":  standing.Name,
				"error": err.Error(),
			}).Error("Failed to upsert standing")
			continue
		}
		successCount++
	}

	log.WithFields(logrus.Fields{
		"success_count": successCount,
		"total_count":   len(standings),
	}).Info("Standings sync completed")

	return nil
}

// SyncSchedules fetches schedules from SportsData.io API and stores them in the database
func (s *Service) SyncSchedules(ctx context.Context, season string) error {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "sportsdata_service.SyncSchedules",
		"season":    season,
	})
	log.Info("Syncing schedules from SportsData.io API to database")

	schedules, err := s.client.GetSchedules(ctx, season)
	if err != nil {
		log.WithError(err).Error("Failed to fetch schedules from API")
		return err
	}

	log.WithField("count", len(schedules)).Info("Upserting schedules in database")

	successCount := 0
	for _, schedule := range schedules {
		_, err := s.schedulesRepo.UpsertByGameKey(ctx, &schedule)
		if err != nil {
			log.WithFields(logrus.Fields{
				"game_key": schedule.GameKey,
				"teams":    schedule.AwayTeam + "@" + schedule.HomeTeam,
				"error":    err.Error(),
			}).Error("Failed to upsert schedule")
			continue
		}
		successCount++
	}

	log.WithFields(logrus.Fields{
		"success_count": successCount,
		"total_count":   len(schedules),
	}).Info("Schedules sync completed")

	return nil
}

// SyncGames fetches games from SportsData.io API and stores them in the database
func (s *Service) SyncGames(ctx context.Context, season string) error {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "sportsdata_service.SyncGames",
		"season":    season,
	})
	log.Info("Syncing games from SportsData.io API to database")

	games, err := s.client.GetGames(ctx, season)
	if err != nil {
		log.WithError(err).Error("Failed to fetch games from API")
		return err
	}

	log.WithField("count", len(games)).Info("Upserting games in database")

	successCount := 0
	for _, game := range games {
		_, err := s.gamesRepo.UpsertByGameKey(ctx, &game)
		if err != nil {
			log.WithFields(logrus.Fields{
				"game_key": game.GameKey,
				"teams":    game.AwayTeam + "@" + game.HomeTeam,
				"error":    err.Error(),
			}).Error("Failed to upsert game")
			continue
		}
		successCount++
	}

	log.WithFields(logrus.Fields{
		"success_count": successCount,
		"total_count":   len(games),
	}).Info("Games sync completed")

	return nil
}

// SyncAll syncs all data for a specified season
func (s *Service) SyncAll(ctx context.Context, season string) error {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "sportsdata_service.SyncAll",
		"season":    season,
	})
	log.Info("Starting complete data sync")

	startTime := time.Now()

	// Sync teams
	if err := s.SyncTeams(ctx); err != nil {
		log.WithError(err).Error("Failed to sync teams")
		return err
	}

	// Sync players
	if err := s.SyncPlayers(ctx); err != nil {
		log.WithError(err).Error("Failed to sync players")
		return err
	}

	// Sync standings
	if err := s.SyncStandings(ctx, season); err != nil {
		log.WithError(err).Error("Failed to sync standings")
		return err
	}

	// Sync schedules
	if err := s.SyncSchedules(ctx, season); err != nil {
		log.WithError(err).Error("Failed to sync schedules")
		return err
	}

	// Sync games
	if err := s.SyncGames(ctx, season); err != nil {
		log.WithError(err).Error("Failed to sync games")
		return err
	}

	duration := time.Since(startTime)
	log.WithField("duration_ms", duration.Milliseconds()).Info("All data synced successfully")

	return nil
}
