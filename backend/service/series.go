package service

import (
	"math"
	"sort"
	"time"

	"dontstarve-stats/alert"
	"dontstarve-stats/cache"
	"dontstarve-stats/model"
)

// snapshots every 15mins * 4 * 24 * 5 = 5 days
const five_days int = 4 * 24 * 5

// GetSeriesContinents returns how many players over time played distributed by continent
func (s Service) GetSeriesContinents() []model.Series {
	if cache.Exists("series_continents") {
		return cache.GetSeries("series_continents")
	}

	snapshots := s.store.GetSeriesContinents()

	m := make(map[string][]model.Item)
	for _, snapshot := range snapshots {
		m["Asia"] = append(m["Asia"], model.Item{snapshot.Date, snapshot.Asia})
		m["Europe"] = append(m["Europe"], model.Item{snapshot.Date, snapshot.Europe})
		m["North America"] = append(m["North America"], model.Item{snapshot.Date, snapshot.NorthAmerica})
		m["South America"] = append(m["South America"], model.Item{snapshot.Date, snapshot.SouthAmerica})
		m["Africa"] = append(m["Africa"], model.Item{snapshot.Date, snapshot.Africa})
		m["Oceania"] = append(m["Oceania"], model.Item{snapshot.Date, snapshot.Oceania})
	}

	result := toSeries(m)[:min(len(m), five_days)]

	cache.SetItems("series_continents", result)

	return result
}

// GetSeriesCharacters returns all the characters and their total players over time
func (s Service) GetSeriesCharacters() []model.Series {
	if cache.Exists("series_characters") {
		return cache.GetSeries("series_characters")
	}

	snapshots := s.store.GetSeriesCharacters()

	m := make(map[string][]model.Item) // character -> [#played at timepoint 0, ... timepoint 1]
	for _, snapshot := range snapshots {
		for _, character := range snapshot.Characters {
			if !isVanillaChar(character.Name) {
				continue
			}
			m[character.Name] = append(m[character.Name], model.Item{snapshot.Date, character.Count})
		}
	}

	result := toSeries(m)[:min(len(m), five_days)]

	cache.SetItems("series_characters", result)

	return result
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
func (s Service) GetCountryCharacters(name string) []model.IsoItem {
	if cache.Exists("country_characters") {
		return cache.GetIso("country_characters")
	}

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

	result := make([]model.IsoItem, 0)
	for country, value := range percentage {
		if total[country] >= 30 {
			result = append(result, model.IsoItem{country, iso[country], value})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i][2].(float64) > result[j][2].(float64)
	})

	final := result[:min(5, len(result))]

	cache.SetItems("country_characters", final)

	return final
}

func (s Service) GetCountryRankings(name string) []model.Series {
	if cache.Exists("country_rankings") {
		return cache.GetSeries("country_rankings")
	}

	ranking := s.store.GetPercentageSnapshot(name)

	// TODO: Handle empty cases
	result := make([]model.Series, 0)
	for _, rank := range ranking {
		data := []model.Item{
			{rank.First, rank.FirstPercent},
			{rank.Second, rank.SecondPercent},
			{rank.Third, rank.ThirdPercent},
			{rank.Fourth, rank.FourthPercent},
			{rank.Fifth, rank.FifthPercent},
		}

		result = append(result, model.Series{
			Name: rank.Date.Format("Jan 2 15:04"),
			Data: data,
		})
	}

	cache.SetItems("country_rankings", result)

	return result
}

func round(f float64) float64 {
	return math.Round(f*100) / 100
}
