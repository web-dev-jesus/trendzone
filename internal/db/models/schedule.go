package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GameKey             string             `bson:"GameKey" json:"gameKey"`
	SeasonType          int                `bson:"SeasonType" json:"seasonType"`
	Season              int                `bson:"Season" json:"season"`
	Week                int                `bson:"Week" json:"week"`
	Date                time.Time          `bson:"Date" json:"date"`
	AwayTeam            string             `bson:"AwayTeam" json:"awayTeam"`
	HomeTeam            string             `bson:"HomeTeam" json:"homeTeam"`
	Channel             string             `bson:"Channel" json:"channel"`
	StadiumID           int                `bson:"StadiumID" json:"stadiumID"`
	Canceled            bool               `bson:"Canceled" json:"canceled"`
	PointSpread         float64            `bson:"PointSpread" json:"pointSpread"`
	OverUnder           float64            `bson:"OverUnder" json:"overUnder"`
	ForecastTempLow     int                `bson:"ForecastTempLow" json:"forecastTempLow"`
	ForecastTempHigh    int                `bson:"ForecastTempHigh" json:"forecastTempHigh"`
	ForecastDescription string             `bson:"ForecastDescription" json:"forecastDescription"`
	ForecastWindSpeed   int                `bson:"ForecastWindSpeed" json:"forecastWindSpeed"`
	AwayTeamMoneyLine   int                `bson:"AwayTeamMoneyLine" json:"awayTeamMoneyLine"`
	HomeTeamMoneyLine   int                `bson:"HomeTeamMoneyLine" json:"homeTeamMoneyLine"`
	Day                 time.Time          `bson:"Day" json:"day"`
	DateTime            time.Time          `bson:"DateTime" json:"dateTime"`
	Status              string             `bson:"Status" json:"status"`
	LastUpdated         time.Time          `bson:"last_updated" json:"lastUpdated"`
}
