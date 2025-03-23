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

const stadiumsCollection = "stadiums"

// StadiumRepository provides methods to interact with stadium data
type StadiumRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewStadiumRepository creates a new stadium repository
func NewStadiumRepository(db *mongo.Database, logger *util.Logger) *StadiumRepository {
	return &StadiumRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertStadium inserts or updates a stadium record
func (r *StadiumRepository) UpsertStadium(ctx context.Context, stadium *models.Stadium) error {
	stadium.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"StadiumID": stadium.StadiumID}
	update := bson.M{"$set": stadium}

	_, err := r.db.Collection(stadiumsCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert stadium", "stadiumID", stadium.StadiumID, "error", err)
		return err
	}

	return nil
}

// GetStadiumByID retrieves a stadium by its ID
func (r *StadiumRepository) GetStadiumByID(ctx context.Context, stadiumID int) (*models.Stadium, error) {
	var stadium models.Stadium
	filter := bson.M{"StadiumID": stadiumID}

	err := r.db.Collection(stadiumsCollection).FindOne(ctx, filter).Decode(&stadium)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get stadium by ID", "stadiumID", stadiumID, "error", err)
		return nil, err
	}

	return &stadium, nil
}

// GetStadiumByName retrieves a stadium by its name
func (r *StadiumRepository) GetStadiumByName(ctx context.Context, name string) (*models.Stadium, error) {
	var stadium models.Stadium
	filter := bson.M{"Name": name}

	err := r.db.Collection(stadiumsCollection).FindOne(ctx, filter).Decode(&stadium)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get stadium by name", "name", name, "error", err)
		return nil, err
	}

	return &stadium, nil
}

// GetAllStadiums retrieves all stadiums
func (r *StadiumRepository) GetAllStadiums(ctx context.Context) ([]*models.Stadium, error) {
	filter := bson.M{}
	cursor, err := r.db.Collection(stadiumsCollection).Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find all stadiums", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var stadiums []*models.Stadium
	if err := cursor.All(ctx, &stadiums); err != nil {
		r.logger.Error("Failed to decode all stadiums", "error", err)
		return nil, err
	}

	return stadiums, nil
}

// DeleteStadium deletes a stadium by ID
func (r *StadiumRepository) DeleteStadium(ctx context.Context, stadiumID int) error {
	filter := bson.M{"StadiumID": stadiumID}
	_, err := r.db.Collection(stadiumsCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete stadium", "stadiumID", stadiumID, "error", err)
		return err
	}

	return nil
}
