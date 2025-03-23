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

const playersCollection = "players"

// PlayerRepository provides methods to interact with player data
type PlayerRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewPlayerRepository creates a new player repository
func NewPlayerRepository(db *mongo.Database, logger *util.Logger) *PlayerRepository {
	return &PlayerRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertPlayer inserts or updates a player record
func (r *PlayerRepository) UpsertPlayer(ctx context.Context, player *models.Player) error {
	player.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"PlayerID": player.PlayerID}
	update := bson.M{"$set": player}

	_, err := r.db.Collection(playersCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert player", "playerID", player.PlayerID, "error", err)
		return err
	}

	return nil
}

// GetPlayerByID retrieves a player by their ID
func (r *PlayerRepository) GetPlayerByID(ctx context.Context, playerID int) (*models.Player, error) {
	var player models.Player
	filter := bson.M{"PlayerID": playerID}

	err := r.db.Collection(playersCollection).FindOne(ctx, filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get player by ID", "playerID", playerID, "error", err)
		return nil, err
	}

	return &player, nil
}

// GetPlayersByTeam retrieves all players for a specific team
func (r *PlayerRepository) GetPlayersByTeam(ctx context.Context, teamKey string) ([]*models.Player, error) {
	filter := bson.M{"Team": teamKey}
	cursor, err := r.db.Collection(playersCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find players by team", "team", teamKey, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var players []*models.Player
	if err := cursor.All(ctx, &players); err != nil {
		r.logger.Error("Failed to decode players by team", "team", teamKey, "error", err)
		return nil, err
	}

	return players, nil
}

// GetPlayersByPosition retrieves all players for a specific position
func (r *PlayerRepository) GetPlayersByPosition(ctx context.Context, position string) ([]*models.Player, error) {
	filter := bson.M{"Position": position}
	cursor, err := r.db.Collection(playersCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find players by position", "position", position, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var players []*models.Player
	if err := cursor.All(ctx, &players); err != nil {
		r.logger.Error("Failed to decode players by position", "position", position, "error", err)
		return nil, err
	}

	return players, nil
}

// GetAllPlayers retrieves all players
func (r *PlayerRepository) GetAllPlayers(ctx context.Context) ([]*models.Player, error) {
	filter := bson.M{}
	cursor, err := r.db.Collection(playersCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find all players", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var players []*models.Player
	if err := cursor.All(ctx, &players); err != nil {
		r.logger.Error("Failed to decode all players", "error", err)
		return nil, err
	}

	return players, nil
}

// DeletePlayer deletes a player by ID
func (r *PlayerRepository) DeletePlayer(ctx context.Context, playerID int) error {
	filter := bson.M{"PlayerID": playerID}
	_, err := r.db.Collection(playersCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete player", "playerID", playerID, "error", err)
		return err
	}

	return nil
}
