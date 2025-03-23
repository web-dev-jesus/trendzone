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

const gamesCollection = "games"

// GameRepository provides methods to interact with game data
type GameRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewGameRepository creates a new game repository
func NewGameRepository(db *mongo.Database, logger *util.Logger) *GameRepository {
	return &GameRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertGame inserts or updates a game record
func (r *GameRepository) UpsertGame(ctx context.Context, game *models.Game) error {
	game.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"GameKey": game.GameKey}
	update := bson.M{"$set": game}

	_, err := r.db.Collection(gamesCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert game", "gameKey", game.GameKey, "error", err)
		return err
	}

	return nil
}

// GetGameByKey retrieves a game by its key
func (r *GameRepository) GetGameByKey(ctx context.Context, gameKey string) (*models.Game, error) {
	var game models.Game
	filter := bson.M{"GameKey": gameKey}

	err := r.db.Collection(gamesCollection).FindOne(ctx, filter).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get game by key", "gameKey", gameKey, "error", err)
		return nil, err
	}

	return &game, nil
}

// GetGamesByWeek retrieves all games for a specific week and season
func (r *GameRepository) GetGamesByWeek(ctx context.Context, season int, seasonType int, week int) ([]*models.Game, error) {
	filter := bson.M{"Season": season, "SeasonType": seasonType, "Week": week}
	cursor, err := r.db.Collection(gamesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find games by week", "season", season, "seasonType", seasonType, "week", week, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var games []*models.Game
	if err := cursor.All(ctx, &games); err != nil {
		r.logger.Error("Failed to decode games by week", "season", season, "seasonType", seasonType, "week", week, "error", err)
		return nil, err
	}

	return games, nil
}

// GetGamesByTeam retrieves all games for a specific team in a season
func (r *GameRepository) GetGamesByTeam(ctx context.Context, team string, season int) ([]*models.Game, error) {
	filter := bson.M{
		"Season": season,
		"$or": []bson.M{
			{"HomeTeam": team},
			{"AwayTeam": team},
		},
	}
	cursor, err := r.db.Collection(gamesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find games by team", "team", team, "season", season, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var games []*models.Game
	if err := cursor.All(ctx, &games); err != nil {
		r.logger.Error("Failed to decode games by team", "team", team, "season", season, "error", err)
		return nil, err
	}

	return games, nil
}

// GetGamesBySeason retrieves all games for a specific season
func (r *GameRepository) GetGamesBySeason(ctx context.Context, season int, seasonType int) ([]*models.Game, error) {
	filter := bson.M{"Season": season, "SeasonType": seasonType}
	cursor, err := r.db.Collection(gamesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find games by season", "season", season, "seasonType", seasonType, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var games []*models.Game
	if err := cursor.All(ctx, &games); err != nil {
		r.logger.Error("Failed to decode games by season", "season", season, "seasonType", seasonType, "error", err)
		return nil, err
	}

	return games, nil
}

// GetLiveGames retrieves all currently in-progress games
func (r *GameRepository) GetLiveGames(ctx context.Context) ([]*models.Game, error) {
	filter := bson.M{"Status": "InProgress"}
	cursor, err := r.db.Collection(gamesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find live games", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var games []*models.Game
	if err := cursor.All(ctx, &games); err != nil {
		r.logger.Error("Failed to decode live games", "error", err)
		return nil, err
	}

	return games, nil
}

// GetGamesByDateRange retrieves all games within a specific date range
func (r *GameRepository) GetGamesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.Game, error) {
	filter := bson.M{
		"Date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	cursor, err := r.db.Collection(gamesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find games by date range", "startDate", startDate, "endDate", endDate, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var games []*models.Game
	if err := cursor.All(ctx, &games); err != nil {
		r.logger.Error("Failed to decode games by date range", "startDate", startDate, "endDate", endDate, "error", err)
		return nil, err
	}

	return games, nil
}

// DeleteGame deletes a game by its key
func (r *GameRepository) DeleteGame(ctx context.Context, gameKey string) error {
	filter := bson.M{"GameKey": gameKey}
	_, err := r.db.Collection(gamesCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete game", "gameKey", gameKey, "error", err)
		return err
	}

	return nil
}
