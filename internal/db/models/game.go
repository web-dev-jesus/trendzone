package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Weather struct {
	Temperature         int    `bson:"Temperature" json:"temperature"`
	Humidity            int    `bson:"Humidity" json:"humidity"`
	WindSpeed           int    `bson:"WindSpeed" json:"windSpeed"`
	ForecastDescription string `bson:"ForecastDescription" json:"forecastDescription"`
}

type Game struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GameKey           string             `bson:"GameKey" json:"gameKey"`
	ScoreID           int                `bson:"ScoreID" json:"scoreID"`
	SeasonType        int                `bson:"SeasonType" json:"seasonType"`
	Season            int                `bson:"Season" json:"season"`
	Week              int                `bson:"Week" json:"week"`
	Date              time.Time          `bson:"Date" json:"date"`
	AwayTeam          string             `bson:"AwayTeam" json:"awayTeam"`
	HomeTeam          string             `bson:"HomeTeam" json:"homeTeam"`
	AwayScore         int                `bson:"AwayScore" json:"awayScore"`
	HomeScore         int                `bson:"HomeScore" json:"homeScore"`
	Channel           string             `bson:"Channel" json:"channel"`
	Stadium           string             `bson:"Stadium" json:"stadium"`
	Status            string             `bson:"Status" json:"status"`
	Quarter           string             `bson:"Quarter" json:"quarter"`
	TimeRemaining     string             `bson:"TimeRemaining" json:"timeRemaining"`
	Possession        *string            `bson:"Possession" json:"possession"`
	Down              *int               `bson:"Down" json:"down"`
	Distance          *int               `bson:"Distance" json:"distance"`
	YardLine          *int               `bson:"YardLine" json:"yardLine"`
	YardLineTerritory *string            `bson:"YardLineTerritory" json:"yardLineTerritory"`
	RedZone           bool               `bson:"RedZone" json:"redZone"`
	AwayTeamMoneyLine int                `bson:"AwayTeamMoneyLine" json:"awayTeamMoneyLine"`
	HomeTeamMoneyLine int                `bson:"HomeTeamMoneyLine" json:"homeTeamMoneyLine"`
	PointSpread       float64            `bson:"PointSpread" json:"pointSpread"`
	OverUnder         float64            `bson:"OverUnder" json:"overUnder"`
	AwayScoreQuarter1 int                `bson:"AwayScoreQuarter1" json:"awayScoreQuarter1"`
	AwayScoreQuarter2 int                `bson:"AwayScoreQuarter2" json:"awayScoreQuarter2"`
	AwayScoreQuarter3 int                `bson:"AwayScoreQuarter3" json:"awayScoreQuarter3"`
	AwayScoreQuarter4 int                `bson:"AwayScoreQuarter4" json:"awayScoreQuarter4"`
	AwayScoreOvertime int                `bson:"AwayScoreOvertime" json:"awayScoreOvertime"`
	HomeScoreQuarter1 int                `bson:"HomeScoreQuarter1" json:"homeScoreQuarter1"`
	HomeScoreQuarter2 int                `bson:"HomeScoreQuarter2" json:"homeScoreQuarter2"`
	HomeScoreQuarter3 int                `bson:"HomeScoreQuarter3" json:"homeScoreQuarter3"`
	HomeScoreQuarter4 int                `bson:"HomeScoreQuarter4" json:"homeScoreQuarter4"`
	HomeScoreOvertime int                `bson:"HomeScoreOvertime" json:"homeScoreOvertime"`
	Weather           Weather            `bson:"Weather" json:"weather"`
	LastUpdated       time.Time          `bson:"last_updated" json:"lastUpdated"`
}
