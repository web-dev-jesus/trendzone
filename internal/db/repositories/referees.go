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

const refereesCollection = "referees"

// RefereeRepository provides methods to interact with referee data
type RefereeRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewRefereeRepository creates a new referee repository
func NewRefereeRepository(db *mongo.Database, logger *util.Logger) *RefereeRepository {
	return &RefereeRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertReferee inserts or updates a referee record
func (r *RefereeRepository) UpsertReferee(ctx context.Context, referee *models.Referee) error {
	referee.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"RefereeID": referee.RefereeID}
	update := bson.M{"$set": referee}

	_, err := r.db.Collection(refereesCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert referee", "refereeID", referee.RefereeID, "error", err)
		return err
	}

	return nil
}

// GetRefereeByID retrieves a referee by their ID
func (r *RefereeRepository) GetRefereeByID(ctx context.Context, refereeID int) (*models.Referee, error) {
	var referee models.Referee
	filter := bson.M{"RefereeID": refereeID}

	err := r.db.Collection(refereesCollection).FindOne(ctx, filter).Decode(&referee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get referee by ID", "refereeID", refereeID, "error", err)
		return nil, err
	}

	return &referee, nil
}

// GetRefereeByName retrieves a referee by their name
func (r *RefereeRepository) GetRefereeByName(ctx context.Context, name string) (*models.Referee, error) {
	var referee models.Referee
	filter := bson.M{"Name": name}

	err := r.db.Collection(refereesCollection).FindOne(ctx, filter).Decode(&referee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get referee by name", "name", name, "error", err)
		return nil, err
	}

	return &referee, nil
}

// GetAllReferees retrieves all referees
func (r *RefereeRepository) GetAllReferees(ctx context.Context) ([]*models.Referee, error) {
	filter := bson.M{}
	cursor, err := r.db.Collection(refereesCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find all referees", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var referees []*models.Referee
	if err := cursor.All(ctx, &referees); err != nil {
		r.logger.Error("Failed to decode all referees", "error", err)
		return nil, err
	}

	return referees, nil
}

// DeleteReferee deletes a referee by ID
func (r *RefereeRepository) DeleteReferee(ctx context.Context, refereeID int) error {
	filter := bson.M{"RefereeID": refereeID}
	_, err := r.db.Collection(refereesCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete referee", "refereeID", refereeID, "error", err)
		return err
	}

	return nil
}
