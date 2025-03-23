package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Standing struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StandingID       int                `bson:"StandingID" json:"standingID"`
	SeasonType       int                `bson:"SeasonType" json:"seasonType"`
	Season           int                `bson:"Season" json:"season"`
	Conference       string             `bson:"Conference" json:"conference"`
	Division         string             `bson:"Division" json:"division"`
	Team             string             `bson:"Team" json:"team"`
	Name             string             `bson:"Name" json:"name"`
	Wins             int                `bson:"Wins" json:"wins"`
	Losses           int                `bson:"Losses" json:"losses"`
	Ties             int                `bson:"Ties" json:"ties"`
	Percentage       float64            `bson:"Percentage" json:"percentage"`
	PointsFor        int                `bson:"PointsFor" json:"pointsFor"`
	PointsAgainst    int                `bson:"PointsAgainst" json:"pointsAgainst"`
	DivisionWins     int                `bson:"DivisionWins" json:"divisionWins"`
	DivisionLosses   int                `bson:"DivisionLosses" json:"divisionLosses"`
	ConferenceWins   int                `bson:"ConferenceWins" json:"conferenceWins"`
	ConferenceLosses int                `bson:"ConferenceLosses" json:"conferenceLosses"`
	DivisionRank     int                `bson:"DivisionRank" json:"divisionRank"`
	ConferenceRank   int                `bson:"ConferenceRank" json:"conferenceRank"`
	HomeWins         int                `bson:"HomeWins" json:"homeWins"`
	HomeLosses       int                `bson:"HomeLosses" json:"homeLosses"`
	AwayWins         int                `bson:"AwayWins" json:"awayWins"`
	AwayLosses       int                `bson:"AwayLosses" json:"awayLosses"`
	LastUpdated      time.Time          `bson:"last_updated" json:"lastUpdated"`
}
