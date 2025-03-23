package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/web-dev-jesus/trendzone/internal/db/models"
	"github.com/web-dev-jesus/trendzone/internal/logger"
)

type PlayersRepository struct {
	collection *mongo.Collection
}

func NewPlayersRepository(client *mongo.Database) *PlayersRepository {
	return &PlayersRepository{
		collection: client.Collection("players"),
	}
}

func (r *PlayersRepository) FindAll(ctx context.Context) ([]models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.FindAll")
	log.Info("Fetching all players")

	var players []models.Player
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.WithError(err).Error("Failed to find players")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &players); err != nil {
		log.WithError(err).Error("Failed to decode players")
		return nil, err
	}

	log.WithField("count", len(players)).Info("Players retrieved successfully")
	return players, nil
}

func (r *PlayersRepository) FindByID(ctx context.Context, id string) (*models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.FindByID").WithField("player_id", id)
	log.Info("Finding player by ID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return nil, err
	}

	var player models.Player
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&player); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Player not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find player")
		return nil, err
	}

	log.Info("Player found")
	return &player, nil
}

func (r *PlayersRepository) FindByTeam(ctx context.Context, teamKey string) ([]models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.FindByTeam").WithField("team_key", teamKey)
	log.Info("Finding players by team")

	var players []models.Player
	cursor, err := r.collection.Find(ctx, bson.M{"Team": teamKey})
	if err != nil {
		log.WithError(err).Error("Failed to find players by team")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &players); err != nil {
		log.WithError(err).Error("Failed to decode players")
		return nil, err
	}

	log.WithField("count", len(players)).Info("Players retrieved successfully")
	return players, nil
}

func (r *PlayersRepository) FindByPlayerID(ctx context.Context, playerID int) (*models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.FindByPlayerID").WithField("player_id", playerID)
	log.Info("Finding player by PlayerID")

	var player models.Player
	if err := r.collection.FindOne(ctx, bson.M{"PlayerID": playerID}).Decode(&player); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Player not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find player")
		return nil, err
	}

	log.Info("Player found")
	return &player, nil
}

func (r *PlayersRepository) Create(ctx context.Context, player *models.Player) (*models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.Create").WithField("player_name", player.Name)
	log.Info("Creating new player")

	player.LastUpdated = time.Now()

	result, err := r.collection.InsertOne(ctx, player)
	if err != nil {
		log.WithError(err).Error("Failed to create player")
		return nil, err
	}

	player.ID = result.InsertedID.(primitive.ObjectID)

	log.WithField("player_id", player.ID.Hex()).Info("Player created successfully")
	return player, nil
}

func (r *PlayersRepository) Update(ctx context.Context, player *models.Player) (*models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.Update").WithField("player_id", player.ID.Hex())
	log.Info("Updating player")

	player.LastUpdated = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": player.ID},
		player,
	)
	if err != nil {
		log.WithError(err).Error("Failed to update player")
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Warn("No player found with given ID")
		return nil, mongo.ErrNoDocuments
	}

	log.Info("Player updated successfully")
	return player, nil
}

func (r *PlayersRepository) UpsertByPlayerID(ctx context.Context, player *models.Player) (*models.Player, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.UpsertByPlayerID").WithField("player_id", player.PlayerID)
	log.Info("Upserting player by PlayerID")

	player.LastUpdated = time.Now()

	filter := bson.M{"PlayerID": player.PlayerID}
	opts := options.Replace().SetUpsert(true)

	result, err := r.collection.ReplaceOne(ctx, filter, player, opts)
	if err != nil {
		log.WithError(err).Error("Failed to upsert player")
		return nil, err
	}

	// If this was a new document (inserted)
	if result.UpsertedID != nil {
		player.ID = result.UpsertedID.(primitive.ObjectID)
		log.WithField("player_id", player.ID.Hex()).Info("Player created successfully")
		return player, nil
	}

	// If this was an existing document (updated)
	var updatedPlayer models.Player
	if err := r.collection.FindOne(ctx, filter).Decode(&updatedPlayer); err != nil {
		log.WithError(err).Error("Failed to retrieve updated player")
		return nil, err
	}

	log.Info("Player updated successfully")
	return &updatedPlayer, nil
}

func (r *PlayersRepository) Delete(ctx context.Context, id string) error {
	log := logger.WithRequestContext(ctx).WithField("component", "players_repository.Delete").WithField("player_id", id)
	log.Info("Deleting player")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.WithError(err).Error("Failed to delete player")
		return err
	}

	if result.DeletedCount == 0 {
		log.Warn("No player found with given ID")
		return mongo.ErrNoDocuments
	}

	log.Info("Player deleted successfully")
	return nil
}
