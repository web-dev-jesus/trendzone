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

const standingsCollection = "standings"

// StandingRepository provides methods to interact with standings data
type StandingRepository struct {
	db     *mongo.Database
	logger *util.Logger
}

// NewStandingRepository creates a new standings repository
func NewStandingRepository(db *mongo.Database, logger *util.Logger) *StandingRepository {
	return &StandingRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertStanding inserts or updates a standing record
func (r *StandingRepository) UpsertStanding(ctx context.Context, standing *models.Standing) error {
	standing.LastUpdated = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"Season":     standing.Season,
		"SeasonType": standing.SeasonType,
		"Team":       standing.Team,
	}
	update := bson.M{"$set": standing}

	_, err := r.db.Collection(standingsCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to upsert standing", "team", standing.Team, "season", standing.Season, "error", err)
		return err
	}

	return nil
}

// UpsertStandings inserts or updates multiple standing records
func (r *StandingRepository) UpsertStandings(ctx context.Context, standings []*models.Standing) error {
	if len(standings) == 0 {
		return nil
	}

	operations := make([]mongo.WriteModel, len(standings))
	for i, standing := range standings {
		standing.LastUpdated = time.Now()
		filter := bson.M{
			"Season":     standing.Season,
			"SeasonType": standing.SeasonType,
			"Team":       standing.Team,
		}
		update := bson.M{"$set": standing}
		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := r.db.Collection(standingsCollection).BulkWrite(ctx, operations, opts)
	if err != nil {
		r.logger.Error("Failed to bulk upsert standings", "count", len(standings), "error", err)
		return err
	}

	return nil
}

// GetStandingByTeam retrieves a standing by team for a specific season
func (r *StandingRepository) GetStandingByTeam(ctx context.Context, team string, season int, seasonType int) (*models.Standing, error) {
	var standing models.Standing
	filter := bson.M{
		"Team":       team,
		"Season":     season,
		"SeasonType": seasonType,
	}

	err := r.db.Collection(standingsCollection).FindOne(ctx, filter).Decode(&standing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get standing by team", "team", team, "season", season, "error", err)
		return nil, err
	}

	return &standing, nil
}

// GetStandingsByDivision retrieves all standings for a specific division and season
func (r *StandingRepository) GetStandingsByDivision(ctx context.Context, conference string, division string, season int, seasonType int) ([]*models.Standing, error) {
	filter := bson.M{
		"Conference": conference,
		"Division":   division,
		"Season":     season,
		"SeasonType": seasonType,
	}
	opts := options.Find().SetSort(bson.M{"DivisionRank": 1})

	cursor, err := r.db.Collection(standingsCollection).Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find standings by division", "conference", conference, "division", division, "season", season, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var standings []*models.Standing
	if err := cursor.All(ctx, &standings); err != nil {
		r.logger.Error("Failed to decode standings by division", "conference", conference, "division", division, "season", season, "error", err)
		return nil, err
	}

	return standings, nil
}

// GetStandingsByConference retrieves all standings for a specific conference and season
func (r *StandingRepository) GetStandingsByConference(ctx context.Context, conference string, season int, seasonType int) ([]*models.Standing, error) {
	filter := bson.M{
		"Conference": conference,
		"Season":     season,
		"SeasonType": seasonType,
	}
	opts := options.Find().SetSort(bson.M{"ConferenceRank": 1})

	cursor, err := r.db.Collection(standingsCollection).Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find standings by conference", "conference", conference, "season", season, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var standings []*models.Standing
	if err := cursor.All(ctx, &standings); err != nil {
		r.logger.Error("Failed to decode standings by conference", "conference", conference, "season", season, "error", err)
		return nil, err
	}

	return standings, nil
}

// GetStandingsBySeason retrieves all standings for a specific season
func (r *StandingRepository) GetStandingsBySeason(ctx context.Context, season int, seasonType int) ([]*models.Standing, error) {
	filter := bson.M{
		"Season":     season,
		"SeasonType": seasonType,
	}
	opts := options.Find().SetSort(bson.M{"Conference": 1, "Division": 1, "DivisionRank": 1})

	cursor, err := r.db.Collection(standingsCollection).Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find standings by season", "season", season, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var standings []*models.Standing
	if err := cursor.All(ctx, &standings); err != nil {
		r.logger.Error("Failed to decode standings by season", "season", season, "error", err)
		return nil, err
	}

	return standings, nil
}

// DeleteStanding deletes a standing by team and season
func (r *StandingRepository) DeleteStanding(ctx context.Context, team string, season int, seasonType int) error {
	filter := bson.M{
		"Team":       team,
		"Season":     season,
		"SeasonType": seasonType,
	}
	_, err := r.db.Collection(standingsCollection).DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete standing", "team", team, "season", season, "error", err)
		return err
	}

	return nil
}
