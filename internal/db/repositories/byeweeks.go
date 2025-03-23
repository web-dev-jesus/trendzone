package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/web-dev-jesus/trendzone/internal/models"
	"github.com/web-dev-jesus/trendzone/internal/util"
)

const byeWeeksCollection = "byeweeks"

// ByeWeekRepository provides methods to interact with bye week data
type ByeWeekRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewByeWeekRepository creates a new bye week repository
func NewByeWeekRepository(db *mongo.Database, logger *util.Logger) *ByeWeekRepository {
	return &ByeWeekRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertByeWeek inserts or updates a bye week record
func (r *ByeWeekRepository) UpsertByeWeek(ctx context.Context, byeWeek *models.ByeWeek) error {
	byeWeek.LastUpdated = time.Now()

	// Generate a ByeID if not present
	if byeWeek.ByeID == "" {
		byeWeek.ByeID = fmt.Sprintf("%d-%d-%s", byeWeek.Season, byeWeek.Week, byeWeek.Team)
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"ByeID": byeWeek.ByeID}
	update := bson.M{"$set": byeWeek}

	_, err := r.db.Collection(byeWeeksCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert bye week", "byeID", byeWeek.ByeID, "error", err)
		return err
	}

	return nil
}

// GetByeWeekByID retrieves a bye week by its ID
func (r *ByeWeekRepository) GetByeWeekByID(ctx context.Context, byeID string) (*models.ByeWeek, error) {
	var byeWeek models.ByeWeek
	filter := bson.M{"ByeID": byeID}

	err := r.db.Collection(byeWeeksCollection).FindOne(ctx, filter).Decode(&byeWeek)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get bye week by ID", "byeID", byeID, "error", err)
		return nil, err
	}

	return &byeWeek, nil
}

// GetByeWeeksForSeason retrieves all bye weeks for a specific season
func (r *ByeWeekRepository) GetByeWeeksForSeason(ctx context.Context, season int) ([]*models.ByeWeek, error) {
	filter := bson.M{"Season": season}
	cursor, err := r.db.Collection(byeWeeksCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find bye weeks for season", "season", season, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var byeWeeks []*models.ByeWeek
	if err := cursor.All(ctx, &byeWeeks); err != nil {
		r.logger.Error("Failed to decode bye weeks for season", "season", season, "error", err)
		return nil, err
	}

	return byeWeeks, nil
}

// GetByeWeeksForWeek retrieves all bye weeks for a specific week in a season
func (r *ByeWeekRepository) GetByeWeeksForWeek(ctx context.Context, season int, week int) ([]*models.ByeWeek, error) {
	filter := bson.M{"Season": season, "Week": week}
	cursor, err := r.db.Collection(byeWeeksCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find bye weeks for week", "season", season, "week", week, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var byeWeeks []*models.ByeWeek
	if err := cursor.All(ctx, &byeWeeks); err != nil {
		r.logger.Error("Failed to decode bye weeks for week", "season", season, "week", week, "error", err)
		return nil, err
	}

	return byeWeeks, nil
}

// GetTeamByeWeek retrieves the bye week for a specific team and season
func (r *ByeWeekRepository) GetTeamByeWeek(ctx context.Context, team string, season int) (*models.ByeWeek, error) {
	var byeWeek models.ByeWeek
	filter := bson.M{"Team": team, "Season": season}

	err := r.db.Collection(byeWeeksCollection).FindOne(ctx, filter).Decode(&byeWeek)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get team bye week", "team", team, "season", season, "error", err)
		return nil, err
	}

	return &byeWeek, nil
}

// DeleteByeWeek deletes a bye week by ID
func (r *ByeWeekRepository) DeleteByeWeek(ctx context.Context, byeID string) error {
	filter := bson.M{"ByeID": byeID}
	_, err := r.db.Collection(byeWeeksCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete bye week", "byeID", byeID, "error", err)
		return err
	}

	return nil
}
