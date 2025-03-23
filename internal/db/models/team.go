package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TeamID               int                `bson:"TeamID" json:"teamID"`
	Key                  string             `bson:"Key" json:"key"`
	City                 string             `bson:"City" json:"city"`
	Name                 string             `bson:"Name" json:"name"`
	Conference           string             `bson:"Conference" json:"conference"`
	Division             string             `bson:"Division" json:"division"`
	FullName             string             `bson:"FullName" json:"fullName"`
	StadiumID            int                `bson:"StadiumID" json:"stadiumID"`
	ByeWeek              int                `bson:"ByeWeek" json:"byeWeek"`
	HeadCoach            string             `bson:"HeadCoach" json:"headCoach"`
	PrimaryColor         string             `bson:"PrimaryColor" json:"primaryColor"`
	SecondaryColor       string             `bson:"SecondaryColor" json:"secondaryColor"`
	TertiaryColor        string             `bson:"TertiaryColor" json:"tertiaryColor"`
	QuaternaryColor      string             `bson:"QuaternaryColor" json:"quaternaryColor"`
	OffensiveCoordinator string             `bson:"OffensiveCoordinator" json:"offensiveCoordinator"`
	DefensiveCoordinator string             `bson:"DefensiveCoordinator" json:"defensiveCoordinator"`
	SpecialTeamsCoach    string             `bson:"SpecialTeamsCoach" json:"specialTeamsCoach"`
	OffensiveScheme      string             `bson:"OffensiveScheme" json:"offensiveScheme"`
	DefensiveScheme      string             `bson:"DefensiveScheme" json:"defensiveScheme"`
	LastUpdated          time.Time          `bson:"last_updated" json:"lastUpdated"`
}
