// trendzone/internal/db/mongodb/repositories/schedules_repo.go
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

type SchedulesRepository struct {
	collection *mongo.Collection
}

func NewSchedulesRepository(client *mongo.Database) *SchedulesRepository {
	return &SchedulesRepository{
		collection: client.Collection("schedules"),
	}
}

func (r *SchedulesRepository) FindAll(ctx context.Context) ([]models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.FindAll")
	log.Info("Fetching all schedules")

	var schedules []models.Schedule
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.WithError(err).Error("Failed to find schedules")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &schedules); err != nil {
		log.WithError(err).Error("Failed to decode schedules")
		return nil, err
	}

	log.WithField("count", len(schedules)).Info("Schedules retrieved successfully")
	return schedules, nil
}

func (r *SchedulesRepository) FindByID(ctx context.Context, id string) (*models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.FindByID").WithField("schedule_id", id)
	log.Info("Finding schedule by ID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return nil, err
	}

	var schedule models.Schedule
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&schedule); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Schedule not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find schedule")
		return nil, err
	}

	log.Info("Schedule found")
	return &schedule, nil
}

func (r *SchedulesRepository) FindByGameKey(ctx context.Context, gameKey string) (*models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.FindByGameKey").WithField("game_key", gameKey)
	log.Info("Finding schedule by GameKey")

	var schedule models.Schedule
	if err := r.collection.FindOne(ctx, bson.M{"GameKey": gameKey}).Decode(&schedule); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Schedule not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find schedule")
		return nil, err
	}

	log.Info("Schedule found")
	return &schedule, nil
}

func (r *SchedulesRepository) FindByTeam(ctx context.Context, team string) ([]models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.FindByTeam").WithField("team", team)
	log.Info("Finding schedules by team")

	filter := bson.M{
		"$or": []bson.M{
			{"HomeTeam": team},
			{"AwayTeam": team},
		},
	}

	var schedules []models.Schedule
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to find schedules by team")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &schedules); err != nil {
		log.WithError(err).Error("Failed to decode schedules")
		return nil, err
	}

	log.WithField("count", len(schedules)).Info("Schedules retrieved successfully")
	return schedules, nil
}

func (r *SchedulesRepository) FindByWeek(ctx context.Context, season int, week int) ([]models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "schedules_repository.FindByWeek",
		"season":    season,
		"week":      week,
	})
	log.Info("Finding schedules by week")

	filter := bson.M{
		"Season": season,
		"Week":   week,
	}

	var schedules []models.Schedule
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to find schedules by week")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &schedules); err != nil {
		log.WithError(err).Error("Failed to decode schedules")
		return nil, err
	}

	log.WithField("count", len(schedules)).Info("Schedules retrieved successfully")
	return schedules, nil
}

func (r *SchedulesRepository) Create(ctx context.Context, schedule *models.Schedule) (*models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.Create").WithField("game_key", schedule.GameKey)
	log.Info("Creating new schedule")

	schedule.LastUpdated = time.Now()

	result, err := r.collection.InsertOne(ctx, schedule)
	if err != nil {
		log.WithError(err).Error("Failed to create schedule")
		return nil, err
	}

	schedule.ID = result.InsertedID.(primitive.ObjectID)

	log.WithField("schedule_id", schedule.ID.Hex()).Info("Schedule created successfully")
	return schedule, nil
}

func (r *SchedulesRepository) Update(ctx context.Context, schedule *models.Schedule) (*models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.Update").WithField("schedule_id", schedule.ID.Hex())
	log.Info("Updating schedule")

	schedule.LastUpdated = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": schedule.ID},
		schedule,
	)
	if err != nil {
		log.WithError(err).Error("Failed to update schedule")
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Warn("No schedule found with given ID")
		return nil, mongo.ErrNoDocuments
	}

	log.Info("Schedule updated successfully")
	return schedule, nil
}

func (r *SchedulesRepository) UpsertByGameKey(ctx context.Context, schedule *models.Schedule) (*models.Schedule, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.UpsertByGameKey").WithField("game_key", schedule.GameKey)
	log.Info("Upserting schedule by GameKey")

	schedule.LastUpdated = time.Now()

	filter := bson.M{"GameKey": schedule.GameKey}
	opts := options.Replace().SetUpsert(true)

	result, err := r.collection.ReplaceOne(ctx, filter, schedule, opts)
	if err != nil {
		log.WithError(err).Error("Failed to upsert schedule")
		return nil, err
	}

	// If this was a new document (inserted)
	if result.UpsertedID != nil {
		schedule.ID = result.UpsertedID.(primitive.ObjectID)
		log.WithField("schedule_id", schedule.ID.Hex()).Info("Schedule created successfully")
		return schedule, nil
	}

	// If this was an existing document (updated)
	var updatedSchedule models.Schedule
	if err := r.collection.FindOne(ctx, filter).Decode(&updatedSchedule); err != nil {
		log.WithError(err).Error("Failed to retrieve updated schedule")
		return nil, err
	}

	log.Info("Schedule updated successfully")
	return &updatedSchedule, nil
}

func (r *SchedulesRepository) Delete(ctx context.Context, id string) error {
	log := logger.WithRequestContext(ctx).WithField("component", "schedules_repository.Delete").WithField("schedule_id", id)
	log.Info("Deleting schedule")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.WithError(err).Error("Failed to delete schedule")
		return err
	}

	if result.DeletedCount == 0 {
		log.Warn("No schedule found with given ID")
		return mongo.ErrNoDocuments
	}

	log.Info("Schedule deleted successfully")
	return nil
}
