package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Team represents NFL team data
type Team struct {
	TeamID               int       `bson:"TeamID" json:"TeamID"`
	Key                  string    `bson:"Key" json:"Key"`
	City                 string    `bson:"City" json:"City"`
	Name                 string    `bson:"Name" json:"Name"`
	Conference           string    `bson:"Conference" json:"Conference"`
	Division             string    `bson:"Division" json:"Division"`
	FullName             string    `bson:"FullName" json:"FullName"`
	StadiumID            int       `bson:"StadiumID" json:"StadiumID"`
	ByeWeek              int       `bson:"ByeWeek" json:"ByeWeek"`
	HeadCoach            string    `bson:"HeadCoach" json:"HeadCoach"`
	PrimaryColor         string    `bson:"PrimaryColor" json:"PrimaryColor"`
	SecondaryColor       string    `bson:"SecondaryColor" json:"SecondaryColor"`
	TertiaryColor        string    `bson:"TertiaryColor" json:"TertiaryColor"`
	QuaternaryColor      string    `bson:"QuaternaryColor" json:"QuaternaryColor"`
	OffensiveCoordinator string    `bson:"OffensiveCoordinator" json:"OffensiveCoordinator"`
	DefensiveCoordinator string    `bson:"DefensiveCoordinator" json:"DefensiveCoordinator"`
	SpecialTeamsCoach    string    `bson:"SpecialTeamsCoach" json:"SpecialTeamsCoach"`
	OffensiveScheme      string    `bson:"OffensiveScheme" json:"OffensiveScheme"`
	DefensiveScheme      string    `bson:"DefensiveScheme" json:"DefensiveScheme"`
	LastUpdated          time.Time `bson:"last_updated" json:"last_updated"`
}

// Player represents NFL player data
type Player struct {
	PlayerID         int       `bson:"PlayerID" json:"PlayerID"`
	Team             string    `bson:"Team" json:"Team"`
	Number           int       `bson:"Number" json:"Number"`
	FirstName        string    `bson:"FirstName" json:"FirstName"`
	LastName         string    `bson:"LastName" json:"LastName"`
	Position         string    `bson:"Position" json:"Position"`
	Status           string    `bson:"Status" json:"Status"`
	Height           string    `bson:"Height" json:"Height"`
	Weight           int       `bson:"Weight" json:"Weight"`
	BirthDate        time.Time `bson:"BirthDate" json:"BirthDate"`
	College          string    `bson:"College" json:"College"`
	Experience       int       `bson:"Experience" json:"Experience"`
	FantasyPosition  string    `bson:"FantasyPosition" json:"FantasyPosition"`
	Active           bool      `bson:"Active" json:"Active"`
	PositionCategory string    `bson:"PositionCategory" json:"PositionCategory"`
	Name             string    `bson:"Name" json:"Name"`
	Age              int       `bson:"Age" json:"Age"`
	PhotoURL         string    `bson:"PhotoURL" json:"PhotoURL"`
	LastUpdated      time.Time `bson:"last_updated" json:"last_updated"`
}

// Weather contains game weather information
type Weather struct {
	Temperature         int    `bson:"Temperature" json:"Temperature"`
	Humidity            int    `bson:"Humidity" json:"Humidity"`
	WindSpeed           int    `bson:"WindSpeed" json:"WindSpeed"`
	ForecastDescription string `bson:"ForecastDescription" json:"ForecastDescription"`
}

// Game represents NFL game data
type Game struct {
	GameKey           string    `bson:"GameKey" json:"GameKey"`
	ScoreID           int       `bson:"ScoreID" json:"ScoreID"`
	SeasonType        int       `bson:"SeasonType" json:"SeasonType"`
	Season            int       `bson:"Season" json:"Season"`
	Week              int       `bson:"Week" json:"Week"`
	Date              time.Time `bson:"Date" json:"Date"`
	AwayTeam          string    `bson:"AwayTeam" json:"AwayTeam"`
	HomeTeam          string    `bson:"HomeTeam" json:"HomeTeam"`
	AwayScore         int       `bson:"AwayScore" json:"AwayScore"`
	HomeScore         int       `bson:"HomeScore" json:"HomeScore"`
	Channel           string    `bson:"Channel" json:"Channel"`
	Stadium           string    `bson:"Stadium" json:"Stadium"`
	Status            string    `bson:"Status" json:"Status"`
	Quarter           string    `bson:"Quarter" json:"Quarter"`
	TimeRemaining     string    `bson:"TimeRemaining" json:"TimeRemaining"`
	Possession        string    `bson:"Possession" json:"Possession"`
	Down              int       `bson:"Down" json:"Down"`
	Distance          int       `bson:"Distance" json:"Distance"`
	YardLine          int       `bson:"YardLine" json:"YardLine"`
	YardLineTerritory string    `bson:"YardLineTerritory" json:"YardLineTerritory"`
	RedZone           bool      `bson:"RedZone" json:"RedZone"`
	AwayTeamMoneyLine int       `bson:"AwayTeamMoneyLine" json:"AwayTeamMoneyLine"`
	HomeTeamMoneyLine int       `bson:"HomeTeamMoneyLine" json:"HomeTeamMoneyLine"`
	PointSpread       float64   `bson:"PointSpread" json:"PointSpread"`
	OverUnder         float64   `bson:"OverUnder" json:"OverUnder"`
	AwayScoreQuarter1 int       `bson:"AwayScoreQuarter1" json:"AwayScoreQuarter1"`
	AwayScoreQuarter2 int       `bson:"AwayScoreQuarter2" json:"AwayScoreQuarter2"`
	AwayScoreQuarter3 int       `bson:"AwayScoreQuarter3" json:"AwayScoreQuarter3"`
	AwayScoreQuarter4 int       `bson:"AwayScoreQuarter4" json:"AwayScoreQuarter4"`
	AwayScoreOvertime int       `bson:"AwayScoreOvertime" json:"AwayScoreOvertime"`
	HomeScoreQuarter1 int       `bson:"HomeScoreQuarter1" json:"HomeScoreQuarter1"`
	HomeScoreQuarter2 int       `bson:"HomeScoreQuarter2" json:"HomeScoreQuarter2"`
	HomeScoreQuarter3 int       `bson:"HomeScoreQuarter3" json:"HomeScoreQuarter3"`
	HomeScoreQuarter4 int       `bson:"HomeScoreQuarter4" json:"HomeScoreQuarter4"`
	HomeScoreOvertime int       `bson:"HomeScoreOvertime" json:"HomeScoreOvertime"`
	Weather           Weather   `bson:"Weather" json:"Weather"`
	LastUpdated       time.Time `bson:"last_updated" json:"last_updated"`
}

// RedZoneStats contains player redzone statistics
type RedZoneStats struct {
	PassingAttempts    int `bson:"PassingAttempts" json:"PassingAttempts"`
	PassingCompletions int `bson:"PassingCompletions" json:"PassingCompletions"`
	PassingYards       int `bson:"PassingYards" json:"PassingYards"`
	PassingTouchdowns  int `bson:"PassingTouchdowns" json:"PassingTouchdowns"`
}

// PlayerGameStats represents player statistics for a specific game
type PlayerGameStats struct {
	PlayerGameID         int          `bson:"PlayerGameID" json:"PlayerGameID"`
	PlayerID             int          `bson:"PlayerID" json:"PlayerID"`
	GameKey              string       `bson:"GameKey" json:"GameKey"`
	SeasonType           int          `bson:"SeasonType" json:"SeasonType"`
	Season               int          `bson:"Season" json:"Season"`
	Week                 int          `bson:"Week" json:"Week"`
	Team                 string       `bson:"Team" json:"Team"`
	Opponent             string       `bson:"Opponent" json:"Opponent"`
	HomeOrAway           string       `bson:"HomeOrAway" json:"HomeOrAway"`
	PassingAttempts      int          `bson:"PassingAttempts" json:"PassingAttempts"`
	PassingCompletions   int          `bson:"PassingCompletions" json:"PassingCompletions"`
	PassingYards         int          `bson:"PassingYards" json:"PassingYards"`
	PassingTouchdowns    int          `bson:"PassingTouchdowns" json:"PassingTouchdowns"`
	PassingInterceptions int          `bson:"PassingInterceptions" json:"PassingInterceptions"`
	PassingSacks         int          `bson:"PassingSacks" json:"PassingSacks"`
	PassingSackYards     int          `bson:"PassingSackYards" json:"PassingSackYards"`
	RushingAttempts      int          `bson:"RushingAttempts" json:"RushingAttempts"`
	RushingYards         int          `bson:"RushingYards" json:"RushingYards"`
	RushingTouchdowns    int          `bson:"RushingTouchdowns" json:"RushingTouchdowns"`
	ReceivingTargets     int          `bson:"ReceivingTargets" json:"ReceivingTargets"`
	Receptions           int          `bson:"Receptions" json:"Receptions"`
	ReceivingYards       int          `bson:"ReceivingYards" json:"ReceivingYards"`
	ReceivingTouchdowns  int          `bson:"ReceivingTouchdowns" json:"ReceivingTouchdowns"`
	Fumbles              int          `bson:"Fumbles" json:"Fumbles"`
	FumblesLost          int          `bson:"FumblesLost" json:"FumblesLost"`
	FantasyPoints        float64      `bson:"FantasyPoints" json:"FantasyPoints"`
	FantasyPointsPPR     float64      `bson:"FantasyPointsPPR" json:"FantasyPointsPPR"`
	RedZoneStats         RedZoneStats `bson:"RedZoneStats" json:"RedZoneStats"`
	LastUpdated          time.Time    `bson:"last_updated" json:"last_updated"`
}

// Standing represents team standings data
type Standing struct {
	StandingID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SeasonType       int                `bson:"SeasonType" json:"SeasonType"`
	Season           int                `bson:"Season" json:"Season"`
	Conference       string             `bson:"Conference" json:"Conference"`
	Division         string             `bson:"Division" json:"Division"`
	Team             string             `bson:"Team" json:"Team"`
	Name             string             `bson:"Name" json:"Name"`
	Wins             int                `bson:"Wins" json:"Wins"`
	Losses           int                `bson:"Losses" json:"Losses"`
	Ties             int                `bson:"Ties" json:"Ties"`
	Percentage       float64            `bson:"Percentage" json:"Percentage"`
	PointsFor        int                `bson:"PointsFor" json:"PointsFor"`
	PointsAgainst    int                `bson:"PointsAgainst" json:"PointsAgainst"`
	DivisionWins     int                `bson:"DivisionWins" json:"DivisionWins"`
	DivisionLosses   int                `bson:"DivisionLosses" json:"DivisionLosses"`
	ConferenceWins   int                `bson:"ConferenceWins" json:"ConferenceWins"`
	ConferenceLosses int                `bson:"ConferenceLosses" json:"ConferenceLosses"`
	DivisionRank     int                `bson:"DivisionRank" json:"DivisionRank"`
	ConferenceRank   int                `bson:"ConferenceRank" json:"ConferenceRank"`
	HomeWins         int                `bson:"HomeWins" json:"HomeWins"`
	HomeLosses       int                `bson:"HomeLosses" json:"HomeLosses"`
	AwayWins         int                `bson:"AwayWins" json:"AwayWins"`
	AwayLosses       int                `bson:"AwayLosses" json:"AwayLosses"`
	LastUpdated      time.Time          `bson:"last_updated" json:"last_updated"`
}

// Schedule represents NFL game schedule information
type Schedule struct {
	GameKey             string    `bson:"GameKey" json:"GameKey"`
	SeasonType          int       `bson:"SeasonType" json:"SeasonType"`
	Season              int       `bson:"Season" json:"Season"`
	Week                int       `bson:"Week" json:"Week"`
	Date                time.Time `bson:"Date" json:"Date"`
	AwayTeam            string    `bson:"AwayTeam" json:"AwayTeam"`
	HomeTeam            string    `bson:"HomeTeam" json:"HomeTeam"`
	Channel             string    `bson:"Channel" json:"Channel"`
	StadiumID           int       `bson:"StadiumID" json:"StadiumID"`
	Canceled            bool      `bson:"Canceled" json:"Canceled"`
	PointSpread         float64   `bson:"PointSpread" json:"PointSpread"`
	OverUnder           float64   `bson:"OverUnder" json:"OverUnder"`
	ForecastTempLow     int       `bson:"ForecastTempLow" json:"ForecastTempLow"`
	ForecastTempHigh    int       `bson:"ForecastTempHigh" json:"ForecastTempHigh"`
	ForecastDescription string    `bson:"ForecastDescription" json:"ForecastDescription"`
	ForecastWindSpeed   int       `bson:"ForecastWindSpeed" json:"ForecastWindSpeed"`
	AwayTeamMoneyLine   int       `bson:"AwayTeamMoneyLine" json:"AwayTeamMoneyLine"`
	HomeTeamMoneyLine   int       `bson:"HomeTeamMoneyLine" json:"HomeTeamMoneyLine"`
	Day                 time.Time `bson:"Day" json:"Day"`
	DateTime            time.Time `bson:"DateTime" json:"DateTime"`
	Status              string    `bson:"Status" json:"Status"`
	LastUpdated         time.Time `bson:"last_updated" json:"last_updated"`
}

// PlayStat represents individual player statistics for a play
type PlayStat struct {
	PlayStatID int    `bson:"PlayStatID" json:"PlayStatID"`
	PlayID     int    `bson:"PlayID" json:"PlayID"`
	Sequence   int    `bson:"Sequence" json:"Sequence"`
	PlayerID   int    `bson:"PlayerID" json:"PlayerID"`
	Name       string `bson:"Name" json:"Name"`
	Team       string `bson:"Team" json:"Team"`
	Opponent   string `bson:"Opponent" json:"Opponent"`
	HomeOrAway string `bson:"HomeOrAway" json:"HomeOrAway"`
	Direction  string `bson:"Direction" json:"Direction"`
	// Include relevant play stat fields based on the data dictionary
}

// Play represents an individual play in a game
type Play struct {
	PlayID               int        `bson:"PlayID" json:"PlayID"`
	QuarterID            int        `bson:"QuarterID" json:"QuarterID"`
	QuarterName          string     `bson:"QuarterName" json:"QuarterName"`
	Sequence             int        `bson:"Sequence" json:"Sequence"`
	TimeRemainingMinutes int        `bson:"TimeRemainingMinutes" json:"TimeRemainingMinutes"`
	TimeRemainingSeconds int        `bson:"TimeRemainingSeconds" json:"TimeRemainingSeconds"`
	PlayTime             time.Time  `bson:"PlayTime" json:"PlayTime"`
	Team                 string     `bson:"Team" json:"Team"`
	Opponent             string     `bson:"Opponent" json:"Opponent"`
	Down                 int        `bson:"Down" json:"Down"`
	Distance             int        `bson:"Distance" json:"Distance"`
	YardLine             int        `bson:"YardLine" json:"YardLine"`
	YardLineTerritory    string     `bson:"YardLineTerritory" json:"YardLineTerritory"`
	YardsToEndZone       int        `bson:"YardsToEndZone" json:"YardsToEndZone"`
	Type                 string     `bson:"Type" json:"Type"`
	YardsGained          int        `bson:"YardsGained" json:"YardsGained"`
	Description          string     `bson:"Description" json:"Description"`
	IsScoringPlay        bool       `bson:"IsScoringPlay" json:"IsScoringPlay"`
	PlayStats            []PlayStat `bson:"PlayStats" json:"PlayStats"`
}

// Quarter represents a quarter in an NFL game
type Quarter struct {
	QuarterID     int    `bson:"QuarterID" json:"QuarterID"`
	ScoreID       int    `bson:"ScoreID" json:"ScoreID"`
	Number        int    `bson:"Number" json:"Number"`
	Name          string `bson:"Name" json:"Name"`
	Description   string `bson:"Description" json:"Description"`
	AwayTeamScore int    `bson:"AwayTeamScore" json:"AwayTeamScore"`
	HomeTeamScore int    `bson:"HomeTeamScore" json:"HomeTeamScore"`
}

// PlayByPlay represents all play-by-play data for a game
type PlayByPlay struct {
	GameKey     string    `bson:"GameKey" json:"GameKey"`
	Score       Game      `bson:"Score" json:"Score"`
	Quarters    []Quarter `bson:"Quarters" json:"Quarters"`
	Plays       []Play    `bson:"Plays" json:"Plays"`
	LastUpdated time.Time `bson:"last_updated" json:"last_updated"`
}

// Stadium represents NFL stadium information
type Stadium struct {
	StadiumID      int       `bson:"StadiumID" json:"StadiumID"`
	Name           string    `bson:"Name" json:"Name"`
	City           string    `bson:"City" json:"City"`
	State          string    `bson:"State" json:"State"`
	Country        string    `bson:"Country" json:"Country"`
	Capacity       int       `bson:"Capacity" json:"Capacity"`
	PlayingSurface string    `bson:"PlayingSurface" json:"PlayingSurface"`
	GeoLat         float64   `bson:"GeoLat" json:"GeoLat"`
	GeoLong        float64   `bson:"GeoLong" json:"GeoLong"`
	Type           string    `bson:"Type" json:"Type"`
	LastUpdated    time.Time `bson:"last_updated" json:"last_updated"`
}

// DepthChartPosition represents a player's position on the depth chart
type DepthChartPosition struct {
	DepthChartID     int    `bson:"DepthChartID" json:"DepthChartID"`
	TeamID           int    `bson:"TeamID" json:"TeamID"`
	PlayerID         int    `bson:"PlayerID" json:"PlayerID"`
	Name             string `bson:"Name" json:"Name"`
	PositionCategory string `bson:"PositionCategory" json:"PositionCategory"`
	Position         string `bson:"Position" json:"Position"`
	DepthOrder       int    `bson:"DepthOrder" json:"DepthOrder"`
}

// DepthChart represents a team's depth chart
type DepthChart struct {
	TeamID       int                  `bson:"TeamID" json:"TeamID"`
	Team         string               `bson:"Team" json:"Team"`
	Offense      []DepthChartPosition `bson:"Offense" json:"Offense"`
	Defense      []DepthChartPosition `bson:"Defense" json:"Defense"`
	SpecialTeams []DepthChartPosition `bson:"SpecialTeams" json:"SpecialTeams"`
	LastUpdated  time.Time            `bson:"last_updated" json:"last_updated"`
}

// Referee represents an NFL referee
type Referee struct {
	RefereeID   int       `bson:"RefereeID" json:"RefereeID"`
	Name        string    `bson:"Name" json:"Name"`
	Number      int       `bson:"Number" json:"Number"`
	Position    string    `bson:"Position" json:"Position"`
	College     string    `bson:"College" json:"College"`
	Experience  int       `bson:"Experience" json:"Experience"`
	LastUpdated time.Time `bson:"last_updated" json:"last_updated"`
}

// ByeWeek represents a team's bye week
type ByeWeek struct {
	ByeID       string    `bson:"ByeID" json:"ByeID"`
	Season      int       `bson:"Season" json:"Season"`
	Week        int       `bson:"Week" json:"Week"`
	Team        string    `bson:"Team" json:"Team"`
	LastUpdated time.Time `bson:"last_updated" json:"last_updated"`
}

// Metadata tracks API calls and updates
type Metadata struct {
	Key       string    `bson:"key" json:"key"`
	Endpoint  string    `bson:"endpoint" json:"endpoint"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Status    string    `bson:"status" json:"status"`
	Notes     string    `bson:"notes" json:"notes"`
}
