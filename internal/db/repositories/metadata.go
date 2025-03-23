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

// MetadataRepository handles database operations for the metadata collection
type MetadataRepository struct {
	collection *mongo.Collection
	logger     *util.Logger
}

// NewMetadataRepository creates a new metadata repository
func NewMetadataRepository(db *mongo.Database, logger *util.Logger) *MetadataRepository {
	return &MetadataRepository{
		collection: db.Collection("metadata"),
		logger:     logger,
	}
}

// RecordAPICall records metadata about an API call
func (r *MetadataRepository) RecordAPICall(key, endpoint, status, notes string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metadata := models.Metadata{
		Key:       key,
		Endpoint:  endpoint,
		Timestamp: time.Now(),
		Status:    status,
		Notes:     notes,
	}

	// Create upsert filter and options
	filter := bson.M{"key": key}
	opts := options.Update().SetUpsert(true)

	// Perform upsert operation
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": metadata}, opts)
	return err
}

// GetLastUpdate retrieves the last update time for an endpoint
func (r *MetadataRepository) GetLastUpdate(key string) (*models.Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var metadata models.Metadata
	err := r.collection.FindOne(ctx, bson.M{"key": key}).Decode(&metadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No metadata found
		}
		return nil, err
	}

	return &metadata, nil
}

// IsUpdateNeeded checks if an update is needed based on hours since last update
func (r *MetadataRepository) IsUpdateNeeded(key string, hours int) (bool, error) {
	metadata, err := r.GetLastUpdate(key)
	if err != nil {
		return false, err
	}

	// If no metadata or timestamp exists, update is needed
	if metadata == nil {
		return true, nil
	}

	// Check if hours have passed since last update
	timeSinceUpdate := time.Since(metadata.Timestamp)
	return timeSinceUpdate.Hours() >= float64(hours), nil
}
