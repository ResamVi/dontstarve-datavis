package model

import "time"

// Server is a (trimmed down in data) representation of a single joinable Don't Starve Together Server
type Server struct {
	ID        int
	Name      string
	Country   string
	Iso       string
	Continent string
	Platform  string
	Mode      string
	Season    string
	Intent    string
	Connected float64
	Mods      bool
	Elapsed   int
	Cycle     int
	Date      time.Time
	Players   []Player
}

// Player is a single connected client belonging to a DST Server
type Player struct {
	ID        int
	Cycle     int
	Name      string
	Character string
	Country   string
	Iso       string
	Continent string
	ServerID  int
}

// ContinentSnapshot keeps a series (of snapshots) of which continent has the most players over time
type ContinentSnapshot struct {
	ID           int
	Date         time.Time
	Asia         int
	Europe       int
	NorthAmerica int
	SouthAmerica int
	Africa       int
	Oceania      int
}

// CharacterSnapshot is keeps a series (of snapshots) of which character is picked most over time
type CharacterSnapshot struct {
	ID         int
	Date       time.Time
	Characters []Character
}

// Character is a (character name, count) pair
type Character struct {
	ID                  int
	Name                string
	Count               int
	CharacterSnapshotID int
}

type PercentageSnapshot struct {
	ID            int
	Date          time.Time
	Character     string
	First         string
	FirstPercent  float64
	Second        string
	SecondPercent float64
	Third         string
	ThirdPercent  float64
	Fourth        string
	FourthPercent float64
	Fifth         string
	FifthPercent  float64
}

// SELECT * FROM character_snapshots INNER JOIN characters ON character_snapshots.id = characters.character_snapshot_id ORDER BY count DESC;

// PlayerStat counts playtime of all characters a character played in seconds
type PlayerStat struct {
	ID   int
	Name string
	//	Playtime []PlayerCharacter
	Wendy        float64
	Wigfrid      float64
	Wilson       float64
	Woodie       float64
	Wolfgang     float64
	Wickerbottom float64
	WX78         float64
	Walter       float64
	Webber       float64
	Winona       float64
	Maxwell      float64
	Wortox       float64
	Wormwood     float64
	Wurt         float64
	Wes          float64
	Willow       float64
	Warly        float64
	Wanda        float64
}

// Keeps a record how long we tracked
type Start struct {
	ID      int
	Started time.Time
}
