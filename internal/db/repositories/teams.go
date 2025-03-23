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

// TeamsRepository handles database operations for the teams collection
type TeamsRepository struct {
	collection *mongo.Collection
	logger     *util.Logger
}

// NewTeamsRepository creates a new teams repository
func NewTeamsRepository(db *mongo.Database, logger *util.Logger) *TeamsRepository {
	return &TeamsRepository{
		collection: db.Collection("teams"),
		logger:     logger,
	}
}

// UpsertTeams inserts or updates teams data
func (r *TeamsRepository) UpsertTeams(teams []models.Team) (int, error) {
	if len(teams) == 0 {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count := 0
	for _, team := range teams {
		// Set last updated timestamp
		team.LastUpdated = time.Now()

		// Create upsert filter and options
		filter := bson.M{"TeamID": team.TeamID}
		opts := options.Update().SetUpsert(true)

		// Perform upsert operation
		result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": team}, opts)
		if err != nil {
			r.logger.Error("Error upserting team: " + err.Error())
			continue
		}

		// Count modified or upserted documents
		if result.ModifiedCount > 0 || result.UpsertedCount > 0 {
			count++
		}
	}

	return count, nil
}

// GetAllTeams retrieves all teams
func (r *TeamsRepository) GetAllTeams() ([]models.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all teams
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var teams []models.Team
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}

// GetTeamByKey retrieves a team by its key (abbreviation)
func (r *TeamsRepository) GetTeamByKey(key string) (*models.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find team by key
	var team models.Team
	err := r.collection.FindOne(ctx, bson.M{"Key": key}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No team found
		}
		return nil, err
	}

	return &team, nil
}

// DeleteTeam deletes a team
func (r *TeamsRepository) DeleteTeam(teamID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete team
	_, err := r.collection.DeleteOne(ctx, bson.M{"TeamID": teamID})
	return err
}
