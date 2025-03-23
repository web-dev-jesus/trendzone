package sportsdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/db/models"
	"github.com/web-dev-jesus/trendzone/internal/logger"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewClient creates a new SportsData.io API client
func NewClient(cfg *config.SportsDataConfig) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
	}
}

// GetTeams retrieves all NFL teams
func (c *Client) GetTeams(ctx context.Context) ([]models.Team, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "sportsdata_client.GetTeams")
	log.Info("Fetching teams from SportsData.io API")

	url := fmt.Sprintf("%s/scores/json/TeamsBasic?key=%s", c.baseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to execute request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithField("status_code", resp.StatusCode).Error("SportsData.io API returned error status")
		return nil, fmt.Errorf("SportsData.io API returned status code %d", resp.StatusCode)
	}

	var teams []models.Team
	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		log.WithError(err).Error("Failed to decode response")
		return nil, err
	}

	// Update LastUpdated for all teams
	now := time.Now()
	for i := range teams {
		teams[i].LastUpdated = now
	}

	log.WithField("count", len(teams)).Info("Successfully fetched teams from API")
	return teams, nil
}

// GetPlayers retrieves all NFL players
func (c *Client) GetPlayers(ctx context.Context) ([]models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "sportsdata_client.GetPlayers")
	log.Info("Fetching players from SportsData.io API")

	url := fmt.Sprintf("%s/scores/json/PlayersByAvailable?key=%s", c.baseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to execute request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithField("status_code", resp.StatusCode).Error("SportsData.io API returned error status")
		return nil, fmt.Errorf("SportsData.io API returned status code %d", resp.StatusCode)
	}

	var players []models.Player
	if err := json.NewDecoder(resp.Body).Decode(&players); err != nil {
		log.WithError(err).Error("Failed to decode response")
		return nil, err
	}

	// Update LastUpdated for all players
	now := time.Now()
	for i := range players {
		players[i].LastUpdated = now
	}

	log.WithField("count", len(players)).Info("Successfully fetched players from API")
	return players, nil
}

// GetStandings retrieves NFL standings for a specified season
func (c *Client) GetStandings(ctx context.Context, season string) ([]models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "sportsdata_client.GetStandings",
		"season":    season,
	})
	log.Info("Fetching standings from SportsData.io API")

	url := fmt.Sprintf("%s/scores/json/Standings/%s?key=%s", c.baseURL, season, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to execute request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithField("status_code", resp.StatusCode).Error("SportsData.io API returned error status")
		return nil, fmt.Errorf("SportsData.io API returned status code %d", resp.StatusCode)
	}

	var standings []models.Standing
	if err := json.NewDecoder(resp.Body).Decode(&standings); err != nil {
		log.WithError(err).Error("Failed to decode response")
		return nil, err
	}

	// Update LastUpdated for all standings
	now := time.Now()
	for i := range standings {
		standings[i].LastUpdated = now
	}

	log.WithField("count", len(standings)).Info("Successfully fetched standings from API")
	return standings, nil
}

// GetSchedules retrieves NFL schedules for a specified season
func (c *Client) GetSchedules(ctx context.Context, season string) ([]models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "sportsdata_client.GetSchedules",
		"season":    season,
	})
	log.Info("Fetching schedules from SportsData.io API")

	url := fmt.Sprintf("%s/scores/json/Schedules/%s?key=%s", c.baseURL, season, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to execute request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithField("status_code", resp.StatusCode).Error("SportsData.io API returned error status")
		return nil, fmt.Errorf("SportsData.io API returned status code %d", resp.StatusCode)
	}

	var schedules []models.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&schedules); err != nil {
		log.WithError(err).Error("Failed to decode response")
		return nil, err
	}

	// Update LastUpdated for all schedules
	now := time.Now()
	for i := range schedules {
		schedules[i].LastUpdated = now
	}

	log.WithField("count", len(schedules)).Info("Successfully fetched schedules from API")
	return schedules, nil
}

// GetGames retrieves NFL games for a specified season (final scores)
func (c *Client) GetGames(ctx context.Context, season string) ([]models.Game, error) {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "sportsdata_client.GetGames",
		"season":    season,
	})
	log.Info("Fetching games from SportsData.io API")

	url := fmt.Sprintf("%s/stats/json/ScoresFinal/%s?key=%s", c.baseURL, season, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to execute request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithField("status_code", resp.StatusCode).Error("SportsData.io API returned error status")
		return nil, fmt.Errorf("SportsData.io API returned status code %d", resp.StatusCode)
	}

	var games []models.Game
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		log.WithError(err).Error("Failed to decode response")
		return nil, err
	}

	// Update LastUpdated for all games
	now := time.Now()
	for i := range games {
		games[i].LastUpdated = now
	}

	log.WithField("count", len(games)).Info("Successfully fetched games from API")
	return games, nil
}
