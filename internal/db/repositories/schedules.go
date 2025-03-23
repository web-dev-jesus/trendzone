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

const schedulesCollection = "schedules"

// ScheduleRepository provides methods to interact with schedule data
type ScheduleRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewScheduleRepository creates a new schedule repository
func NewScheduleRepository(db *mongo.Database, logger *util.Logger) *ScheduleRepository {
	return &ScheduleRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertSchedule inserts or updates a schedule record
func (r *ScheduleRepository) UpsertSchedule(ctx context.Context, schedule *models.Schedule) error {
	schedule.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"GameKey": schedule.GameKey}
	update := bson.M{"$set": schedule}

	_, err := r.db.Collection(schedulesCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert schedule", "gameKey", schedule.GameKey, "error", err)
		return err
	}

	return nil
}

// UpsertSchedules inserts or updates multiple schedule records
func (r *ScheduleRepository) UpsertSchedules(ctx context.Context, schedules []*models.Schedule) error {
	if len(schedules) == 0 {
		return nil
	}

	operations := make([]mongo.WriteModel, len(schedules))
	for i, schedule := range schedules {
		schedule.LastUpdated = time.Now()
		filter := bson.M{"GameKey": schedule.GameKey}
		update := bson.M{"$set": schedule}
		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.db.Collection(schedulesCollection).BulkWrite(ctx, operations, opts)
	if err != nil {
		r.logger.Error("Failed to bulk upsert schedules", "count", len(schedules), "error", err)
		return err
	}

	return nil
}

// GetScheduleByGameKey retrieves a schedule by its game key
func (r *ScheduleRepository) GetScheduleByGameKey(ctx context.Context, gameKey string) (*models.Schedule, error) {
	var schedule models.Schedule
	filter := bson.M{"GameKey": gameKey}

	err := r.db.Collection(schedulesCollection).FindOne(ctx, filter).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get schedule by game key", "gameKey", gameKey, "error", err)
		return nil, err
	}

	return &schedule, nil
}

// GetSchedulesByWeek retrieves all schedules for a specific week and season
func (r *ScheduleRepository) GetSchedulesByWeek(ctx context.Context, season int, seasonType int, week int) ([]*models.Schedule, error) {
	filter := bson.M{"Season": season, "SeasonType": seasonType, "Week": week}
	cursor, err := r.db.Collection(schedulesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find schedules by week", "season", season, "seasonType", seasonType, "week", week, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*models.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		r.logger.Error("Failed to decode schedules by week", "season", season, "seasonType", seasonType, "week", week, "error", err)
		return nil, err
	}

	return schedules, nil
}

// GetSchedulesByTeam retrieves all schedules for a specific team in a season
func (r *ScheduleRepository) GetSchedulesByTeam(ctx context.Context, team string, season int) ([]*models.Schedule, error) {
	filter := bson.M{
		"Season": season,
		"$or": []bson.M{
			{"HomeTeam": team},
			{"AwayTeam": team},
		},
	}
	cursor, err := r.db.Collection(schedulesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find schedules by team", "team", team, "season", season, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*models.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		r.logger.Error("Failed to decode schedules by team", "team", team, "season", season, "error", err)
		return nil, err
	}

	return schedules, nil
}

// GetSchedulesBySeason retrieves all schedules for a specific season
func (r *ScheduleRepository) GetSchedulesBySeason(ctx context.Context, season int, seasonType int) ([]*models.Schedule, error) {
	filter := bson.M{"Season": season, "SeasonType": seasonType}
	cursor, err := r.db.Collection(schedulesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find schedules by season", "season", season, "seasonType", seasonType, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*models.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		r.logger.Error("Failed to decode schedules by season", "season", season, "seasonType", seasonType, "error", err)
		return nil, err
	}

	return schedules, nil
}

// GetUpcomingGames retrieves upcoming games starting from a specific date
func (r *ScheduleRepository) GetUpcomingGames(ctx context.Context, startDate time.Time, limit int) ([]*models.Schedule, error) {
	filter := bson.M{"Date": bson.M{"$gte": startDate}}
	opts := options.Find().SetSort(bson.M{"Date": 1}).SetLimit(int64(limit))

	cursor, err := r.db.Collection(schedulesCollection).Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find upcoming games", "startDate", startDate, "limit", limit, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*models.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		r.logger.Error("Failed to decode upcoming games", "startDate", startDate, "limit", limit, "error", err)
		return nil, err
	}

	return schedules, nil
}

// DeleteSchedule deletes a schedule by its game key
func (r *ScheduleRepository) DeleteSchedule(ctx context.Context, gameKey string) error {
	filter := bson.M{"GameKey": gameKey}
	_, err := r.db.Collection(schedulesCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete schedule", "gameKey", gameKey, "error", err)
		return err
	}

	return nil
}
