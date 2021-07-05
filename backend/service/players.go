package service

import (
	"time"

	"dontstarve-stats/cache"
	"dontstarve-stats/model"
)

// CountPlayers returns the total players online
func (s Service) CountPlayers() int {
	if cache.Exists("player_count") {
		val, err := cache.Get("player_count").Int()
		if err != nil {
			panic(err)
		}

		return val
	}

	players := s.store.GetAllPlayers()

	cache.Set("player_count", len(players))

	return len(players)
}

// CountCharacters returns how often each character is picked
func (s Service) CountCharacters(includeModded bool) []model.Item {
	if cache.Exists("character_count") && !includeModded {
		return cache.GetItems("character_count")
	}

	players := s.store.GetAllPlayers()

	m := make(map[string]int)
	for _, player := range players {
		if includeModded || isVanillaChar(player.Character) {
			m[player.Character]++
		}
	}
	result := toItems(m)

	if !includeModded {
		cache.SetItems("character_count", result)
	}

	return result
}

// CountPlayerOrigin returns the countries where the majority of players come from (the top 20 highest countries only)
// Note: Player inherit their origin from the server's origin (i.e. its IP), so German players on french servers are French
func (s Service) CountPlayerOrigin(includeAll bool) []model.Item {
	if cache.Exists("player_origin") && !includeAll {
		return cache.GetItems("player_origin")
	}

	players := s.store.GetAllPlayers()

	m := make(map[string]int)
	for _, player := range players {
		m[player.Country]++
	}

	if includeAll {
		return toItems(m)
	}

	result := toItems(m)[:min(20, len(m))]

	cache.SetItems("player_origin", result)

	return result
}

// GetCountryPreference returns how often each character is picked for a specific country
func (s Service) GetCountryPreference(country string) []model.Item {
	if cache.Exists(country) {
		return cache.GetItems(country)
	}

	players := s.store.GetAllPlayers()

	m := make(map[string]int)
	for _, player := range players {
		if player.Country != country {
			continue
		}

		if !isVanillaChar(player.Character) {
			continue
		}

		m[player.Character]++
	}

	result := toItems(m)

	cache.SetItems(country, result)

	return result
}

// Given the list of current players compare to the previous
// list (assumes we havent updated the table of players yet - run before svc.ServerSnapshot) and if we see players again add to their playtime
func (s Service) TrackPlaytime(servers []model.Server, lastCheck time.Time) {
	previous := s.store.GetAllPlayers()

	prev := make(map[string]model.Player)
	for _, player := range previous {
		prev[player.Name] = player
	}

	for _, server := range servers {
		for _, player := range server.Players {
			if _, exists := prev[player.Name]; exists {

				stat := s.store.GetPlayerStat(player.Name)
				// TODO: UGLY
				switch player.Character {
				case "Wendy":
					stat.Wendy += time.Since(lastCheck).Seconds()
				case "Wigfrid":
					stat.Wigfrid += time.Since(lastCheck).Seconds()
				case "Wilson":
					stat.Wilson += time.Since(lastCheck).Seconds()
				case "Woodie":
					stat.Woodie += time.Since(lastCheck).Seconds()
				case "Wolfgang":
					stat.Wolfgang += time.Since(lastCheck).Seconds()
				case "Wickerbottom":
					stat.Wickerbottom += time.Since(lastCheck).Seconds()
				case "WX-78":
					stat.WX78 += time.Since(lastCheck).Seconds()
				case "Walter":
					stat.Walter += time.Since(lastCheck).Seconds()
				case "Webber":
					stat.Webber += time.Since(lastCheck).Seconds()
				case "Winona":
					stat.Winona += time.Since(lastCheck).Seconds()
				case "Maxwell":
					stat.Maxwell += time.Since(lastCheck).Seconds()
				case "Wortox":
					stat.Wortox += time.Since(lastCheck).Seconds()
				case "Wormwood":
					stat.Wormwood += time.Since(lastCheck).Seconds()
				case "Wurt":
					stat.Wurt += time.Since(lastCheck).Seconds()
				case "Wes":
					stat.Wes += time.Since(lastCheck).Seconds()
				case "Willow":
					stat.Willow += time.Since(lastCheck).Seconds()
				case "Warly":
					stat.Warly += time.Since(lastCheck).Seconds()
				default:
					return
				}

				s.store.UpdatePlayerStat(stat)
			}
		}

	}
}

// GetPlayTime returns how much a players plays each character
func (s Service) GetPlayTime(name string) []model.Item {

	stats := s.store.GetPlayerStat(name)

	toHours := float64(60 * 60)
	return sortFloatDescending([]model.Item{
		{"Wendy", round(stats.Wendy / toHours)},
		{"Wigfrid", round(stats.Wigfrid / toHours)},
		{"Wilson", round(stats.Wilson / toHours)},
		{"Woodie", round(stats.Woodie / toHours)},
		{"Wolfgang", round(stats.Wolfgang / toHours)},
		{"Wickerbottom", round(stats.Wickerbottom / toHours)},
		{"WX-78", round(stats.WX78 / toHours)},
		{"Walter", round(stats.Walter / toHours)},
		{"Webber", round(stats.Webber / toHours)},
		{"Winona", round(stats.Winona / toHours)},
		{"Maxwell", round(stats.Maxwell / toHours)},
		{"Wortox", round(stats.Wortox / toHours)},
		{"Wormwood", round(stats.Wormwood / toHours)},
		{"Wurt", round(stats.Wurt / toHours)},
		{"Wes", round(stats.Wes / toHours)},
		{"Willow", round(stats.Willow / toHours)},
		{"Warly", round(stats.Warly / toHours)},
	})
}

// dst specific domain knowledge
func isVanillaChar(name string) bool {
	switch name {
	case "Wendy":
		fallthrough
	case "Wigfrid":
		fallthrough
	case "Wilson":
		fallthrough
	case "Woodie":
		fallthrough
	case "Wolfgang":
		fallthrough
	case "Wickerbottom":
		fallthrough
	case "WX-78":
		fallthrough
	case "Walter":
		fallthrough
	case "Webber":
		fallthrough
	case "Winona":
		fallthrough
	case "Maxwell":
		fallthrough
	case "Wortox":
		fallthrough
	case "Wormwood":
		fallthrough
	case "Wurt":
		fallthrough
	case "Wes":
		fallthrough
	case "Willow":
		fallthrough
	case "Warly":
		return true
	default:
		return false
	}
}
