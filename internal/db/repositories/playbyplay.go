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

const playByPlayCollection = "playbyplay"

// PlayByPlayRepository provides methods to interact with play-by-play data
type PlayByPlayRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewPlayByPlayRepository creates a new play-by-play repository
func NewPlayByPlayRepository(db *mongo.Database, logger *util.Logger) *PlayByPlayRepository {
	return &PlayByPlayRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertPlayByPlay inserts or updates a play-by-play record
func (r *PlayByPlayRepository) UpsertPlayByPlay(ctx context.Context, playByPlay *models.PlayByPlay) error {
	playByPlay.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"GameKey": playByPlay.GameKey}
	update := bson.M{"$set": playByPlay}

	_, err := r.db.Collection(playByPlayCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert play-by-play", "gameKey", playByPlay.GameKey, "error", err)
		return err
	}

	return nil
}

// GetPlayByPlayByGameKey retrieves play-by-play data for a specific game
func (r *PlayByPlayRepository) GetPlayByPlayByGameKey(ctx context.Context, gameKey string) (*models.PlayByPlay, error) {
	var playByPlay models.PlayByPlay
	filter := bson.M{"GameKey": gameKey}

	err := r.db.Collection(playByPlayCollection).FindOne(ctx, filter).Decode(&playByPlay)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get play-by-play by game key", "gameKey", gameKey, "error", err)
		return nil, err
	}

	return &playByPlay, nil
}

// GetPlaysByPlayerID retrieves all plays involving a specific player
func (r *PlayByPlayRepository) GetPlaysByPlayerID(ctx context.Context, playerID int) ([]models.Play, error) {
	filter := bson.M{"Plays.PlayStats.PlayerID": playerID}
	cursor, err := r.db.Collection(playByPlayCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find plays by player ID", "playerID", playerID, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var playByPlays []models.PlayByPlay
	if err := cursor.All(ctx, &playByPlays); err != nil {
		r.logger.Error("Failed to decode plays by player ID", "playerID", playerID, "error", err)
		return nil, err
	}

	var plays []models.Play
	for _, pbp := range playByPlays {
		for _, play := range pbp.Plays {
			for _, stat := range play.PlayStats {
				if stat.PlayerID == playerID {
					plays = append(plays, play)
					break
				}
			}
		}
	}

	return plays, nil
}

// GetScoringPlaysByGameKey retrieves all scoring plays for a specific game
func (r *PlayByPlayRepository) GetScoringPlaysByGameKey(ctx context.Context, gameKey string) ([]models.Play, error) {
	filter := bson.M{"GameKey": gameKey, "Plays.IsScoringPlay": true}
	var playByPlay models.PlayByPlay

	err := r.db.Collection(playByPlayCollection).FindOne(ctx, filter).Decode(&playByPlay)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Play{}, nil
		}
		r.logger.Error("Failed to get scoring plays by game key", "gameKey", gameKey, "error", err)
		return nil, err
	}

	var scoringPlays []models.Play
	for _, play := range playByPlay.Plays {
		if play.IsScoringPlay {
			scoringPlays = append(scoringPlays, play)
		}
	}

	return scoringPlays, nil
}

// DeletePlayByPlay deletes play-by-play data for a specific game
func (r *PlayByPlayRepository) DeletePlayByPlay(ctx context.Context, gameKey string) error {
	filter := bson.M{"GameKey": gameKey}
	_, err := r.db.Collection(playByPlayCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete play-by-play", "gameKey", gameKey, "error", err)
		return err
	}

	return nil
}
