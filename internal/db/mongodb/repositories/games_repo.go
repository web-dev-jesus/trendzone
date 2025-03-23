package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/web-dev-jesus/trendzone/internal/db/models"
	"github.com/web-dev-jesus/trendzone/internal/logger"
)

type GamesRepository struct {
	collection *mongo.Collection
}

func NewGamesRepository(client *mongo.Database) *GamesRepository {
	return &GamesRepository{
		collection: client.Collection("games"),
	}
}

func (r *GamesRepository) FindAll(ctx context.Context) ([]models.Game, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.FindAll")
	log.Info("Fetching all games")

	var games []models.Game
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.WithError(err).Error("Failed to find games")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &games); err != nil {
		log.WithError(err).Error("Failed to decode games")
		return nil, err
	}

	log.WithField("count", len(games)).Info("Games retrieved successfully")
	return games, nil
}

func (r *GamesRepository) FindByID(ctx context.Context, id string) (*models.Game, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.FindByID").WithField("game_id", id)
	log.Info("Finding game by ID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return nil, err
	}

	var game models.Game
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&game); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Game not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find game")
		return nil, err
	}

	log.Info("Game found")
	return &game, nil
}

func (r *GamesRepository) FindByGameKey(ctx context.Context, gameKey string) (*models.Game, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.FindByGameKey").WithField("game_key", gameKey)
	log.Info("Finding game by GameKey")

	var game models.Game
	if err := r.collection.FindOne(ctx, bson.M{"GameKey": gameKey}).Decode(&game); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Game not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find game")
		return nil, err
	}

	log.Info("Game found")
	return &game, nil
}

func (r *GamesRepository) FindByTeam(ctx context.Context, team string) ([]models.Game, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.FindByTeam").WithField("team", team)
	log.Info("Finding games by team")

	filter := bson.M{
		"$or": []bson.M{
			{"HomeTeam": team},
			{"AwayTeam": team},
		},
	}

	var games []models.Game
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to find games by team")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &games); err != nil {
		log.WithError(err).Error("Failed to decode games")
		return nil, err
	}

	log.WithField("count", len(games)).Info("Games retrieved successfully")
	return games, nil
}

func (r *GamesRepository) FindByWeek(ctx context.Context, season int, week int) ([]models.Game, error) {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "games_repository.FindByWeek",
		"season":    season,
		"week":      week,
	})
	log.Info("Finding games by week")

	filter := bson.M{
		"Season": season,
		"Week":   week,
	}

	var games []models.Game
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to find games by week")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &games); err != nil {
		log.WithError(err).Error("Failed to decode games")
		return nil, err
	}

	log.WithField("count", len(games)).Info("Games retrieved successfully")
	return games, nil
}

func (r *GamesRepository) Create(ctx context.Context, game *models.Game) (*models.Game, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.Create").WithField("game_key", game.GameKey)
	log.Info("Creating new game")

	game.LastUpdated = time.Now()

	result, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		log.WithError(err).Error("Failed to create game")
		return nil, err
	}

	game.ID = result.InsertedID.(primitive.ObjectID)

	log.WithField("game_id", game.ID.Hex()).Info("Game created successfully")
	return game, nil
}

func (r *GamesRepository) Update(ctx context.Context, game *models.Game) (*models.Game, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.Update").WithField("game_id", game.ID.Hex())
	log.Info("Updating game")

	game.LastUpdated = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": game.ID},
		game,
	)
	if err != nil {
		log.WithError(err).Error("Failed to update game")
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Warn("No game found with given ID")
		return nil, mongo.ErrNoDocuments
	}

	log.Info("Game updated successfully")
	return game, nil
}

func (r *GamesRepository) UpsertByGameKey(ctx context.Context, game *models.Game) (*models.Game, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.UpsertByGameKey").WithField("game_key", game.GameKey)
	log.Info("Upserting game by GameKey")

	game.LastUpdated = time.Now()

	filter := bson.M{"GameKey": game.GameKey}
	opts := options.Replace().SetUpsert(true)

	result, err := r.collection.ReplaceOne(ctx, filter, game, opts)
	if err != nil {
		log.WithError(err).Error("Failed to upsert game")
		return nil, err
	}

	// If this was a new document (inserted)
	if result.UpsertedID != nil {
		game.ID = result.UpsertedID.(primitive.ObjectID)
		log.WithField("game_id", game.ID.Hex()).Info("Game created successfully")
		return game, nil
	}

	// If this was an existing document (updated)
	var updatedGame models.Game
	if err := r.collection.FindOne(ctx, filter).Decode(&updatedGame); err != nil {
		log.WithError(err).Error("Failed to retrieve updated game")
		return nil, err
	}

	log.Info("Game updated successfully")
	return &updatedGame, nil
}

func (r *GamesRepository) Delete(ctx context.Context, id string) error {
	log := logger.WithRequestContext(ctx).WithField("component", "games_repository.Delete").WithField("game_id", id)
	log.Info("Deleting game")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.WithError(err).Error("Failed to delete game")
		return err
	}

	if result.DeletedCount == 0 {
		log.Warn("No game found with given ID")
		return mongo.ErrNoDocuments
	}

	log.Info("Game deleted successfully")
	return nil
}
