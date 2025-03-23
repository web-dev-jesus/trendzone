package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PlayerID         int                `bson:"PlayerID" json:"playerID"`
	Team             string             `bson:"Team" json:"team"`
	Number           int                `bson:"Number" json:"number"`
	FirstName        string             `bson:"FirstName" json:"firstName"`
	LastName         string             `bson:"LastName" json:"lastName"`
	Position         string             `bson:"Position" json:"position"`
	Status           string             `bson:"Status" json:"status"`
	Height           string             `bson:"Height" json:"height"`
	Weight           int                `bson:"Weight" json:"weight"`
	BirthDate        time.Time          `bson:"BirthDate" json:"birthDate"`
	College          string             `bson:"College" json:"college"`
	Experience       int                `bson:"Experience" json:"experience"`
	FantasyPosition  string             `bson:"FantasyPosition" json:"fantasyPosition"`
	Active           bool               `bson:"Active" json:"active"`
	PositionCategory string             `bson:"PositionCategory" json:"positionCategory"`
	Name             string             `bson:"Name" json:"name"`
	Age              int                `bson:"Age" json:"age"`
	PhotoUrl         string             `bson:"PhotoUrl" json:"photoUrl,omitempty"`
	LastUpdated      time.Time          `bson:"last_updated" json:"lastUpdated"`
}
