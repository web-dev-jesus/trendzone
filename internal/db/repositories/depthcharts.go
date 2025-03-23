package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/web-dev-jesus/trendzone/internal/models"
	"github.com/web-dev-jesus/trendzone/internal/util"
)

const depthChartsCollection = "depthcharts"

// DepthChartRepository provides methods to interact with depth chart data
type DepthChartRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewDepthChartRepository creates a new depth chart repository
func NewDepthChartRepository(db *mongo.Database, logger *util.Logger) *DepthChartRepository {
	return &DepthChartRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertDepthChart inserts or updates a depth chart record
func (r *DepthChartRepository) UpsertDepthChart(ctx context.Context, depthChart *models.DepthChart) error {
	depthChart.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"TeamID": depthChart.TeamID}
	update := bson.M{"$set": depthChart}

	_, err := r.db.Collection(depthChartsCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert depth chart", "teamID", depthChart.TeamID, "error", err)
		return err
	}

	return nil
}

// GetDepthChartByTeamID retrieves a depth chart by team ID
func (r *DepthChartRepository) GetDepthChartByTeamID(ctx context.Context, teamID int) (*models.DepthChart, error) {
	var depthChart models.DepthChart
	filter := bson.M{"TeamID": teamID}

	err := r.db.Collection(depthChartsCollection).FindOne(ctx, filter).Decode(&depthChart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get depth chart by team ID", "teamID", teamID, "error", err)
		return nil, err
	}

	return &depthChart, nil
}

// GetDepthChartByTeamAbbreviation retrieves a depth chart by team abbreviation
func (r *DepthChartRepository) GetDepthChartByTeamAbbreviation(ctx context.Context, team string) (*models.DepthChart, error) {
	var depthChart models.DepthChart
	filter := bson.M{"Team": team}

	err := r.db.Collection(depthChartsCollection).FindOne(ctx, filter).Decode(&depthChart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get depth chart by team abbreviation", "team", team, "error", err)
		return nil, err
	}

	return &depthChart, nil
}

// GetAllDepthCharts retrieves all depth charts
func (r *DepthChartRepository) GetAllDepthCharts(ctx context.Context) ([]*models.DepthChart, error) {
	filter := bson.M{}
	cursor, err := r.db.Collection(depthChartsCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find all depth charts", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var depthCharts []*models.DepthChart
	if err := cursor.All(ctx, &depthCharts); err != nil {
		r.logger.Error("Failed to decode all depth charts", "error", err)
		return nil, err
	}

	return depthCharts, nil
}

// GetPlayerDepthChartStatus retrieves a player's position and depth order across all teams
func (r *DepthChartRepository) GetPlayerDepthChartStatus(ctx context.Context, playerID int) ([]map[string]interface{}, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"Offense.PlayerID": playerID},
			{"Defense.PlayerID": playerID},
			{"SpecialTeams.PlayerID": playerID},
		},
	}

	cursor, err := r.db.Collection(depthChartsCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find depth charts for player", "playerID", playerID, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var depthCharts []*models.DepthChart
	if err := cursor.All(ctx, &depthCharts); err != nil {
		r.logger.Error("Failed to decode depth charts for player", "playerID", playerID, "error", err)
		return nil, err
	}

	var result []map[string]interface{}

	for _, dc := range depthCharts {
		// Check offense positions
		for _, pos := range dc.Offense {
			if pos.PlayerID == playerID {
				result = append(result, map[string]interface{}{
					"Team":             dc.Team,
					"Category":         "Offense",
					"Position":         pos.Position,
					"DepthOrder":       pos.DepthOrder,
					"PositionCategory": pos.PositionCategory,
				})
			}
		}

		// Check defense positions
		for _, pos := range dc.Defense {
			if pos.PlayerID == playerID {
				result = append(result, map[string]interface{}{
					"Team":             dc.Team,
					"Category":         "Defense",
					"Position":         pos.Position,
					"DepthOrder":       pos.DepthOrder,
					"PositionCategory": pos.PositionCategory,
				})
			}
		}

		// Check special teams positions
		for _, pos := range dc.SpecialTeams {
			if pos.PlayerID == playerID {
				result = append(result, map[string]interface{}{
					"Team":             dc.Team,
					"Category":         "SpecialTeams",
					"Position":         pos.Position,
					"DepthOrder":       pos.DepthOrder,
					"PositionCategory": pos.PositionCategory,
				})
			}
		}
	}

	return result, nil
}

// DeleteDepthChart deletes a depth chart by team ID
func (r *DepthChartRepository) DeleteDepthChart(ctx context.Context, teamID int) error {
	filter := bson.M{"TeamID": teamID}
	_, err := r.db.Collection(depthChartsCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete depth chart", "teamID", teamID, "error", err)
		return err
	}

	return nil
}
