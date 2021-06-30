package service

import (
	"math"
	"sort"
	"time"

	"dontstarve-stats/alert"
	"dontstarve-stats/model"
)

// SeriesRepository documents all required methods to a database
// necessary to serve series data
type SeriesRepository interface {
}

// GetSeriesContinents returns how many players over time played distributed by continent
func (s Service) GetSeriesContinents() []Series {
	snapshots := s.store.GetSeriesContinents()

	m := make(map[string][]Item)
	for _, snapshot := range snapshots {
		m["Asia"] = append(m["Asia"], Item{snapshot.Date, snapshot.Asia})
		m["Europe"] = append(m["Europe"], Item{snapshot.Date, snapshot.Europe})
		m["North America"] = append(m["North America"], Item{snapshot.Date, snapshot.NorthAmerica})
		m["South America"] = append(m["South America"], Item{snapshot.Date, snapshot.SouthAmerica})
		m["Africa"] = append(m["Africa"], Item{snapshot.Date, snapshot.Africa})
		m["Oceania"] = append(m["Oceania"], Item{snapshot.Date, snapshot.Oceania})
	}

	return toSeries(m)
}

// GetSeriesCharacters returns all the characters and their total players over time
func (s Service) GetSeriesCharacters() []Series {
	snapshots := s.store.GetSeriesCharacters()

	m := make(map[string][]Item) // character -> [#played at timepoint 0, ... timepoint 1]
	for _, snapshot := range snapshots {
		for _, character := range snapshot.Characters {
			if !isVanillaChar(character.Name) {
				continue
			}
			m[character.Name] = append(m[character.Name], Item{snapshot.Date, character.Count})
		}
	}

	return toSeries(m)
}

// ContinentSnapshot creates an entry in a timeseries table of
// how the players are distributed across the globe at the current time
func (s Service) ContinentSnapshot() {
	players := s.store.GetAllPlayers()

	m := make(map[string]int) // continent -> #total players in continent
	for _, player := range players {
		m[player.Continent]++
	}

	s.store.CreateContinentSnapshot(m)
}

// CharacterSnapshot creates an entry in a timeseries table of
// how many characters are played at the current ime
func (s Service) CharacterSnapshot() {
	players := s.store.GetAllPlayers()

	m := make(map[string]int) // character -> #count played
	for _, player := range players {
		m[player.Character]++
	}

	s.store.CreateCharacterSnapshot(m)
}

func (s Service) PercentageSnapshot() {
	CHARACTERS := []string{
		"Wilson", "Willow", "Wolfgang",
		"Wendy", "WX-78", "Wickerbottom",
		"Woodie", "Wes", "Maxwell",
		"Wigfrid", "Webber", "Warly",
		"Wormwood", "Winona", "Wortox",
		"Wurt", "Walter",
	}

	for _, character := range CHARACTERS {
		ranking := s.GetCountryCharacters(character)

		if len(ranking) < 5 {
			alert.Msg("PercentageSnapshot() unsuccesful: too few in ranking")
			return
		}

		s.store.CreatePercentageSnapshot(model.PercentageSnapshot{
			Date:          time.Now().Local(),
			Character:     character,
			First:         ranking[0][1].(string),
			FirstPercent:  ranking[0][2].(float64),
			Second:        ranking[1][1].(string),
			SecondPercent: ranking[1][2].(float64),
			Third:         ranking[2][1].(string),
			ThirdPercent:  ranking[2][2].(float64),
			Fourth:        ranking[3][1].(string),
			FourthPercent: ranking[3][2].(float64),
			Fifth:         ranking[4][1].(string),
			FifthPercent:  ranking[4][2].(float64),
		})
	}
}

//
func (s Service) GetCountryCharacters(name string) []IsoItem {
	players := s.store.GetAllPlayers()

	count := make(map[string]int)  // country -> count of character (`name`) played in this country
	total := make(map[string]int)  // country -> total players in this country
	iso := make(map[string]string) // country -> iso of country

	for _, player := range players {
		count[player.Country] = 0
		total[player.Country] = 0
		iso[player.Country] = player.Iso
	}

	for _, player := range players {
		if player.Character == name {
			count[player.Country]++
		}

		total[player.Country]++
	}

	percentage := make(map[string]float64)
	for country := range count {
		if total[country] == 0 {
			percentage[country] = 0
			continue
		}
		f := float64(count[country]) / float64(total[country])
		percentage[country] = round(f * 100)
	}

	result := make([]IsoItem, 0)
	for country, value := range percentage {
		if total[country] >= 30 {
			result = append(result, IsoItem{country, iso[country], value})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i][2].(float64) > result[j][2].(float64)
	})

	return result[:min(5, len(result))]
}

func (s Service) GetCountryRankings(name string) []Series {
	ranking := s.store.GetPercentageSnapshot(name)

	// TODO: Handle empty cases
	result := make([]Series, 0)
	for _, rank := range ranking {
		data := []Item{
			{rank.First, rank.FirstPercent},
			{rank.Second, rank.SecondPercent},
			{rank.Third, rank.ThirdPercent},
			{rank.Fourth, rank.FourthPercent},
			{rank.Fifth, rank.FifthPercent},
		}

		result = append(result, Series{
			Name: rank.Date.Format("Jan 2 15:04"),
			Data: data,
		})
	}

	return result
}

func round(f float64) float64 {
	return math.Round(f*100) / 100
}
