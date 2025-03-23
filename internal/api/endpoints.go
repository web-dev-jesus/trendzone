// internal/api/endpoints.go

package api

import (
	"fmt"
)

// Endpoint paths for SportsData.io NFL API
const (
	// Base endpoints
	TeamsEndpoint    = "/scores/json/TeamsBasic"
	PlayersEndpoint  = "/scores/json/PlayersByAvailable"
	StadiumsEndpoint = "/scores/json/Stadiums"
	RefereesEndpoint = "/scores/json/Referees"

	// Season-specific endpoints
	StandingsEndpoint = "/scores/json/Standings/%s"
	SchedulesEndpoint = "/scores/json/Schedules/%s"
	ByeWeeksEndpoint  = "/scores/json/Byes/%s"

	// Week-specific endpoints
	ScoresEndpoint      = "/stats/json/ScoresFinal/%s/%d"
	PlayerStatsEndpoint = "/stats/json/PlayerGameStatsByTeamFinal/%s/%d/%s"
	PlayByPlayEndpoint  = "/pbp/json/PlayByPlayFinal/%s/%d/%s"

	// Other endpoints
	DepthChartsEndpoint = "/scores/json/DepthCharts"
)

// GetStandingsEndpoint returns the formatted standings endpoint for a season
func GetStandingsEndpoint(seasonParam string) string {
	return fmt.Sprintf(StandingsEndpoint, seasonParam)
}

// GetSchedulesEndpoint returns the formatted schedules endpoint for a season
func GetSchedulesEndpoint(seasonParam string) string {
	return fmt.Sprintf(SchedulesEndpoint, seasonParam)
}

// GetByeWeeksEndpoint returns the formatted bye weeks endpoint for a season
func GetByeWeeksEndpoint(seasonParam string) string {
	return fmt.Sprintf(ByeWeeksEndpoint, seasonParam)
}

// GetScoresEndpoint returns the formatted scores endpoint for a season and week
func GetScoresEndpoint(seasonParam string, week int) string {
	return fmt.Sprintf(ScoresEndpoint, seasonParam, week)
}

// GetPlayerStatsEndpoint returns the formatted player stats endpoint for a season, week, and team
func GetPlayerStatsEndpoint(seasonParam string, week int, team string) string {
	return fmt.Sprintf(PlayerStatsEndpoint, seasonParam, week, team)
}

// GetPlayByPlayEndpoint returns the formatted play by play endpoint for a season, week, and home team
func GetPlayByPlayEndpoint(seasonParam string, week int, homeTeam string) string {
	return fmt.Sprintf(PlayByPlayEndpoint, seasonParam, week, homeTeam)
}
