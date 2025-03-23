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

const playerGameStatsCollection = "playergamestats"

// PlayerGameStatsRepository provides methods to interact with player game statistics
type PlayerGameStatsRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewPlayerGameStatsRepository creates a new player game stats repository
func NewPlayerGameStatsRepository(db *mongo.Database, logger *util.Logger) *PlayerGameStatsRepository {
	return &PlayerGameStatsRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertPlayerGameStats inserts or updates player game stats
func (r *PlayerGameStatsRepository) UpsertPlayerGameStats(ctx context.Context, stats *models.PlayerGameStats) error {
	stats.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"PlayerGameID": stats.PlayerGameID}
	update := bson.M{"$set": stats}

	_, err := r.db.Collection(playerGameStatsCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert player game stats", "playerGameID", stats.PlayerGameID, "error", err)
		return err
	}

	return nil
}

// UpsertMultiplePlayerGameStats inserts or updates multiple player game stats
func (r *PlayerGameStatsRepository) UpsertMultiplePlayerGameStats(ctx context.Context, statsArray []*models.PlayerGameStats) error {
	if len(statsArray) == 0 {
		return nil
	}

	operations := make([]mongo.WriteModel, len(statsArray))
	for i, stats := range statsArray {
		stats.LastUpdated = time.Now()
		filter := bson.M{"PlayerGameID": stats.PlayerGameID}
		update := bson.M{"$set": stats}
		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.db.Collection(playerGameStatsCollection).BulkWrite(ctx, operations, opts)
	if err != nil {
		r.logger.Error("Failed to bulk upsert player game stats", "count", len(statsArray), "error", err)
		return err
	}

	return nil
}

// GetPlayerGameStatsByID retrieves player game stats by ID
func (r *PlayerGameStatsRepository) GetPlayerGameStatsByID(ctx context.Context, playerGameID int) (*models.PlayerGameStats, error) {
	var stats models.PlayerGameStats
	filter := bson.M{"PlayerGameID": playerGameID}

	err := r.db.Collection(playerGameStatsCollection).FindOne(ctx, filter).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get player game stats by ID", "playerGameID", playerGameID, "error", err)
		return nil, err
	}

	return &stats, nil
}

// GetPlayerGameStatsByPlayerAndGame retrieves player game stats by player ID and game key
func (r *PlayerGameStatsRepository) GetPlayerGameStatsByPlayerAndGame(ctx context.Context, playerID int, gameKey string) (*models.PlayerGameStats, error) {
	var stats models.PlayerGameStats
	filter := bson.M{"PlayerID": playerID, "GameKey": gameKey}

	err := r.db.Collection(playerGameStatsCollection).FindOne(ctx, filter).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get player game stats by player and game", "playerID", playerID, "gameKey", gameKey, "error", err)
		return nil, err
	}

	return &stats, nil
}

// GetPlayerGameStatsByGame retrieves all player game stats for a specific game
func (r *PlayerGameStatsRepository) GetPlayerGameStatsByGame(ctx context.Context, gameKey string) ([]*models.PlayerGameStats, error) {
	filter := bson.M{"GameKey": gameKey}
	cursor, err := r.db.Collection(playerGameStatsCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find player game stats by game", "gameKey", gameKey, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var statsArray []*models.PlayerGameStats
	if err := cursor.All(ctx, &statsArray); err != nil {
		r.logger.Error("Failed to decode player game stats by game", "gameKey", gameKey, "error", err)
		return nil, err
	}

	return statsArray, nil
}

// GetPlayerGameStatsByPlayer retrieves all game stats for a specific player in a season
func (r *PlayerGameStatsRepository) GetPlayerGameStatsByPlayer(ctx context.Context, playerID int, season int) ([]*models.PlayerGameStats, error) {
	filter := bson.M{"PlayerID": playerID, "Season": season}
	opts := options.Find().SetSort(bson.M{"Week": 1})

	cursor, err := r.db.Collection(playerGameStatsCollection).Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find player game stats by player", "playerID", playerID, "season", season, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var statsArray []*models.PlayerGameStats
	if err := cursor.All(ctx, &statsArray); err != nil {
		r.logger.Error("Failed to decode player game stats by player", "playerID", playerID, "season", season, "error", err)
		return nil, err
	}

	return statsArray, nil
}

// GetPlayerGameStatsByTeamAndWeek retrieves all player game stats for a team in a specific week
func (r *PlayerGameStatsRepository) GetPlayerGameStatsByTeamAndWeek(ctx context.Context, team string, season int, week int) ([]*models.PlayerGameStats, error) {
	filter := bson.M{"Team": team, "Season": season, "Week": week}
	cursor, err := r.db.Collection(playerGameStatsCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find player game stats by team and week", "team", team, "season", season, "week", week, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var statsArray []*models.PlayerGameStats
	if err := cursor.All(ctx, &statsArray); err != nil {
		r.logger.Error("Failed to decode player game stats by team and week", "team", team, "season", season, "week", week, "error", err)
		return nil, err
	}

	return statsArray, nil
}

// GetTopPerformersForWeek retrieves top performers for a specific statistic in a week
func (r *PlayerGameStatsRepository) GetTopPerformersForWeek(ctx context.Context, season int, week int, statField string, limit int) ([]*models.PlayerGameStats, error) {
	filter := bson.M{"Season": season, "Week": week}
	opts := options.Find().SetSort(bson.M{statField: -1}).SetLimit(int64(limit))

	cursor, err := r.db.Collection(playerGameStatsCollection).Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find top performers for week", "season", season, "week", week, "statField", statField, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var statsArray []*models.PlayerGameStats
	if err := cursor.All(ctx, &statsArray); err != nil {
		r.logger.Error("Failed to decode top performers for week", "season", season, "week", week, "statField", statField, "error", err)
		return nil, err
	}

	return statsArray, nil
}

// DeletePlayerGameStats deletes player game stats by ID
func (r *PlayerGameStatsRepository) DeletePlayerGameStats(ctx context.Context, playerGameID int) error {
	filter := bson.M{"PlayerGameID": playerGameID}
	_, err := r.db.Collection(playerGameStatsCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete player game stats", "playerGameID", playerGameID, "error", err)
		return err
	}

	return nil
}
