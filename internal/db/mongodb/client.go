// nfl-stats-service/internal/db/mongodb/client.go
package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/web-dev-jesus/trendzone/config"
	"github.com/web-dev-jesus/trendzone/internal/logger"
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewClient establishes a connection to MongoDB
func NewClient(ctx context.Context, cfg *config.MongoDBConfig) (*Client, error) {
	log := logger.WithRequestContext(ctx).WithField("component", "mongodb.client")
	log.Info("Connecting to MongoDB")

	// Create a new client and connect to the server
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		log.WithError(err).Error("Failed to connect to MongoDB")
		return nil, err
	}

	// Ping the primary to verify that the connection is working
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.WithError(err).Error("Failed to ping MongoDB")
		return nil, err
	}

	log.Info("Successfully connected to MongoDB")
	return &Client{
		client: client,
		db:     client.Database(cfg.DBName),
	}, nil
}

// Close disconnects from MongoDB
func (c *Client) Close(ctx context.Context) error {
	log := logger.WithRequestContext(ctx).WithField("component", "mongodb.client")
	log.Info("Closing MongoDB connection")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := c.client.Disconnect(ctx); err != nil {
		log.WithError(err).Error("Failed to close MongoDB connection")
		return err
	}

	log.Info("MongoDB connection closed")
	return nil
}

// GetDatabase returns the database instance
func (c *Client) GetDatabase() *mongo.Database {
	return c.db
}

// GetCollection returns a collection from the database
func (c *Client) GetCollection(name string) *mongo.Collection {
	return c.db.Collection(name)
}
