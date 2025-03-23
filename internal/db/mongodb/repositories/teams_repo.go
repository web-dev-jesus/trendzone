// trendzone/internal/db/mongodb/repositories/teams_repo.go
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

type TeamsRepository struct {
	collection *mongo.Collection
}

func NewTeamsRepository(client *mongo.Database) *TeamsRepository {
	return &TeamsRepository{
		collection: client.Collection("teams"),
	}
}

func (r *TeamsRepository) FindAll(ctx context.Context) ([]models.Team, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "teams_repository.FindAll")
	log.Info("Fetching all teams")

	var teams []models.Team
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.WithError(err).Error("Failed to find teams")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &teams); err != nil {
		log.WithError(err).Error("Failed to decode teams")
		return nil, err
	}

	log.WithField("count", len(teams)).Info("Teams retrieved successfully")
	return teams, nil
}

func (r *TeamsRepository) FindByID(ctx context.Context, id string) (*models.Team, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "teams_repository.FindByID").WithField("team_id", id)
	log.Info("Finding team by ID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return nil, err
	}

	var team models.Team
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&team); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Team not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find team")
		return nil, err
	}

	log.Info("Team found")
	return &team, nil
}

func (r *TeamsRepository) FindByKey(ctx context.Context, key string) (*models.Team, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "teams_repository.FindByKey").WithField("team_key", key)
	log.Info("Finding team by key")

	var team models.Team
	if err := r.collection.FindOne(ctx, bson.M{"Key": key}).Decode(&team); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info("Team not found")
			return nil, nil
		}
		log.WithError(err).Error("Failed to find team")
		return nil, err
	}

	log.Info("Team found")
	return &team, nil
}

func (r *TeamsRepository) Create(ctx context.Context, team *models.Team) (*models.Team, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "teams_repository.Create").WithField("team_key", team.Key)
	log.Info("Creating new team")

	team.LastUpdated = time.Now()

	result, err := r.collection.InsertOne(ctx, team)
	if err != nil {
		log.WithError(err).Error("Failed to create team")
		return nil, err
	}

	team.ID = result.InsertedID.(primitive.ObjectID)

	log.WithField("team_id", team.ID.Hex()).Info("Team created successfully")
	return team, nil
}

func (r *TeamsRepository) Update(ctx context.Context, team *models.Team) (*models.Team, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "teams_repository.Update").WithField("team_id", team.ID.Hex())
	log.Info("Updating team")

	team.LastUpdated = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": team.ID},
		team,
	)
	if err != nil {
		log.WithError(err).Error("Failed to update team")
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Warn("No team found with given ID")
		return nil, mongo.ErrNoDocuments
	}

	log.Info("Team updated successfully")
	return team, nil
}

func (r *TeamsRepository) UpsertByTeamID(ctx context.Context, team *models.Team) (*models.Team, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "teams_repository.UpsertByTeamID").WithField("team_id", team.TeamID)
	log.Info("Upserting team by TeamID")

	team.LastUpdated = time.Now()

	filter := bson.M{"TeamID": team.TeamID}
	opts := options.Replace().SetUpsert(true)

	result, err := r.collection.ReplaceOne(ctx, filter, team, opts)
	if err != nil {
		log.WithError(err).Error("Failed to upsert team")
		return nil, err
	}

	// If this was a new document (inserted)
	if result.UpsertedID != nil {
		team.ID = result.UpsertedID.(primitive.ObjectID)
		log.WithField("team_id", team.ID.Hex()).Info("Team created successfully")
		return team, nil
	}

	// If this was an existing document (updated)
	var updatedTeam models.Team
	if err := r.collection.FindOne(ctx, filter).Decode(&updatedTeam); err != nil {
		log.WithError(err).Error("Failed to retrieve updated team")
		return nil, err
	}

	log.Info("Team updated successfully")
	return &updatedTeam, nil
}

func (r *TeamsRepository) Delete(ctx context.Context, id string) error {
	log := logger.WithRequestContext(ctx).WithField("component", "teams_repository.Delete").WithField("team_id", id)
	log.Info("Deleting team")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Error("Invalid ObjectID format")
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.WithError(err).Error("Failed to delete team")
		return err
	}

	if result.DeletedCount == 0 {
		log.Warn("No team found with given ID")
		return mongo.ErrNoDocuments
	}

	log.Info("Team deleted successfully")
	return nil
}
