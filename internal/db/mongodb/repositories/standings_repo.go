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

type StandingsRepository struct {
	collection *mongo.Collection
}

func NewStandingsRepository(client *mongo.Database) *StandingsRepository {
	return &StandingsRepository{
		collection: client.Collection("standings"),
	}
}

func (r *StandingsRepository) FindAll(ctx context.Context) ([]models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "standings_repository.FindAll")
	log.Info("Fetching all standings")

	var standings []models.Standing
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.WithError(err).Error("Failed to find standings")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &standings); err != nil {
		log.WithError(err).Error("Failed to decode standings")
		return nil, err
	}

	log.WithField("count", len(standings)).Info("Standings retrieved successfully")
	return standings, nil
}

func (r *StandingsRepository) FindByID(ctx context.Context, id string) (*models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "standings_repository.FindByID").WithField("standing_id", id)
	log.Info("Finding standing by ID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return nil, err
	}

	var standing models.Standing
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&standing); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Standing not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find standing")
		return nil, err
	}

	log.Info("Standing found")
	return &standing, nil
}

func (r *StandingsRepository) FindByTeam(ctx context.Context, team string) (*models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "standings_repository.FindByTeam").WithField("team", team)
	log.Info("Finding standing by team")

	var standing models.Standing
	if err := r.collection.FindOne(ctx, bson.M{"Team": team}).Decode(&standing); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Standing not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find standing")
		return nil, err
	}

	log.Info("Standing found")
	return &standing, nil
}

func (r *StandingsRepository) FindByDivision(ctx context.Context, conference string, division string) ([]models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component":  "standings_repository.FindByDivision",
		"conference": conference,
		"division":   division,
	})
	log.Info("Finding standings by division")

	filter := bson.M{
		"Conference": conference,
		"Division":   division,
	}

	var standings []models.Standing
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to find standings by division")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &standings); err != nil {
		log.WithError(err).Error("Failed to decode standings")
		return nil, err
	}

	log.WithField("count", len(standings)).Info("Standings retrieved successfully")
	return standings, nil
}

func (r *StandingsRepository) Create(ctx context.Context, standing *models.Standing) (*models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "standings_repository.Create").WithField("team", standing.Team)
	log.Info("Creating new standing")

	standing.LastUpdated = time.Now()

	result, err := r.collection.InsertOne(ctx, standing)
	if err != nil {
		log.WithError(err).Error("Failed to create standing")
		return nil, err
	}

	standing.ID = result.InsertedID.(primitive.ObjectID)

	log.WithField("standing_id", standing.ID.Hex()).Info("Standing created successfully")
	return standing, nil
}

func (r *StandingsRepository) Update(ctx context.Context, standing *models.Standing) (*models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "standings_repository.Update").WithField("standing_id", standing.ID.Hex())
	log.Info("Updating standing")

	standing.LastUpdated = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": standing.ID},
		standing,
	)
	if err != nil {
		log.WithError(err).Error("Failed to update standing")
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Warn("No standing found with given ID")
		return nil, mongo.ErrNoDocuments
	}

	log.Info("Standing updated successfully")
	return standing, nil
}

func (r *StandingsRepository) UpsertByTeamAndSeason(ctx context.Context, standing *models.Standing) (*models.Standing, error) {
	log := logger.WithRequestContext(ctx).WithFields(logrus.Fields{
		"component": "standings_repository.UpsertByTeamAndSeason",
		"team":      standing.Team,
		"season":    standing.Season,
	})
	log.Info("Upserting standing by Team and Season")

	standing.LastUpdated = time.Now()

	filter := bson.M{
		"Team":   standing.Team,
		"Season": standing.Season,
	}
	opts := options.Replace().SetUpsert(true)

	result, err := r.collection.ReplaceOne(ctx, filter, standing, opts)
	if err != nil {
		log.WithError(err).Error("Failed to upsert standing")
		return nil, err
	}

	// If this was a new document (inserted)
	if result.UpsertedID != nil {
		standing.ID = result.UpsertedID.(primitive.ObjectID)
		log.WithField("standing_id", standing.ID.Hex()).Info("Standing created successfully")
		return standing, nil
	}

	// If this was an existing document (updated)
	var updatedStanding models.Standing
	if err := r.collection.FindOne(ctx, filter).Decode(&updatedStanding); err != nil {
		log.WithError(err).Error("Failed to retrieve updated standing")
		return nil, err
	}

	log.Info("Standing updated successfully")
	return &updatedStanding, nil
}

func (r *StandingsRepository) Delete(ctx context.Context, id string) error {
	log := logger.WithRequestContext(ctx).WithField("component", "standings_repository.Delete").WithField("standing_id", id)
	log.Info("Deleting standing")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.WithError(err).Error("Failed to delete standing")
		return err
	}

	if result.DeletedCount == 0 {
		log.Warn("No standing found with given ID")
		return mongo.ErrNoDocuments
	}

	log.Info("Standing deleted successfully")
	return nil
}
