package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/util"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	clientOnce sync.Once
)

// Connect establishes a connection to MongoDB
func Connect(cfg *config.Config, logger *util.Logger) (*mongo.Database, error) {
	var connectErr error

	clientOnce.Do(func() {
		logger.Info("Connecting to MongoDB...")

		// Set client options with security best practices
		clientOptions := options.Client().
			ApplyURI(cfg.MongoURI).
			SetRetryWrites(true).
			SetRetryReads(true).
			SetConnectTimeout(10 * time.Second).
			SetServerSelectionTimeout(5 * time.Second).
			SetMaxConnIdleTime(60 * time.Second).
			SetMaxPoolSize(100)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, connectErr = mongo.Connect(ctx, clientOptions)
		if connectErr != nil {
			connectErr = fmt.Errorf("failed to connect to MongoDB: %w", connectErr)
			return
		}

		// Ping the database to verify connection
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			connectErr = fmt.Errorf("failed to ping MongoDB: %w", err)
			return
		}

		// Get database instance
		database = client.Database(cfg.MongoDatabase)

		logger.Info(fmt.Sprintf("Connected to MongoDB database: %s", cfg.MongoDatabase))

		// Create indexes
		if err := CreateIndexes(client, cfg.MongoDatabase); err != nil {
			logger.Warn(fmt.Sprintf("Failed to create indexes: %v", err))
		}
	})

	if connectErr != nil {
		return nil, connectErr
	}

	return database, nil
}

// Disconnect closes the MongoDB connection
func Disconnect(logger *util.Logger) error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Disconnecting from MongoDB...")
	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	logger.Info("Disconnected from MongoDB")
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(name string) *mongo.Collection {
	if database == nil {
		return nil
	}
	return database.Collection(name)
}
