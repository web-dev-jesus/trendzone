package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates all required indexes in MongoDB
func CreateIndexes(client *mongo.Client, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := client.Database(dbName)

	// Teams Collection
	teamsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "TeamID", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "Key", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "Conference", Value: 1},
				{Key: "Division", Value: 1},
			},
		},
	}
	if _, err := db.Collection("teams").Indexes().CreateMany(ctx, teamsIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for teams collection")

	// Players Collection
	playersIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "PlayerID", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "Team", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "Position", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "Name", Value: 1}},
		},
	}
	if _, err := db.Collection("players").Indexes().CreateMany(ctx, playersIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for players collection")

	// Games Collection
	gamesIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "GameKey", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "SeasonType", Value: 1},
				{Key: "Season", Value: 1},
				{Key: "Week", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "HomeTeam", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "AwayTeam", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "Date", Value: 1}},
		},
	}
	if _, err := db.Collection("games").Indexes().CreateMany(ctx, gamesIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for games collection")

	// Player Game Stats Collection
	playerGameStatsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "PlayerGameID", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "PlayerID", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "GameKey", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "Team", Value: 1},
				{Key: "Season", Value: 1},
				{Key: "Week", Value: 1},
			},
		},
	}
	if _, err := db.Collection("player_game_stats").Indexes().CreateMany(ctx, playerGameStatsIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for player_game_stats collection")

	// Standings Collection
	standingsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "Season", Value: 1},
				{Key: "SeasonType", Value: 1},
				{Key: "Team", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "Conference", Value: 1},
				{Key: "Division", Value: 1},
			},
		},
	}
	if _, err := db.Collection("standings").Indexes().CreateMany(ctx, standingsIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for standings collection")

	// Schedules Collection
	schedulesIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "GameKey", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "Season", Value: 1},
				{Key: "Week", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "Date", Value: 1}},
		},
	}
	if _, err := db.Collection("schedules").Indexes().CreateMany(ctx, schedulesIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for schedules collection")

	// Play By Play Collection
	playByPlayIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "GameKey", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "Plays.PlayStats.PlayerID", Value: 1}},
		},
	}
	if _, err := db.Collection("play_by_play").Indexes().CreateMany(ctx, playByPlayIndexes); err != nil {
		return err
	}
	log.Println("Created indexes for play_by_play collection")

	// Create indexes for other collections
	// ... (similar pattern for remaining collections)

	return nil
}
