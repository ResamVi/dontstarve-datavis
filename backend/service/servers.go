package service

import (
	"math"
	"time"

	"dontstarve-stats/cache"
	"dontstarve-stats/model"
)

func (s Service) LastUpdate() float64 {
	servers := s.store.GetAllServers()

	if len(servers) == 0 {
		return 0
	}

	return math.Round(time.Since(servers[0].Date).Minutes())
}

// CountServers returns the total servers that are listed
func (s Service) CountServers() int {
	if cache.Exists("server_count") {
		val, err := cache.Get("server_count").Int()
		if err != nil {
			panic(err)
		}

		return val
	}

	servers := s.store.GetAllServers()

	cache.Set("server_count", len(servers))

	return len(servers)
}

// ServerSnapshot updates the list of current public players and servers (= a snapshot)
func (s Service) ServerSnapshot(servers []model.Server) {
	s.store.DeleteAllPlayers()
	s.store.DeleteAllServers()
	s.store.CreateServers(servers)
}

// CountCountry returns how many servers belong to the country (the top 20 highest only)
func (s Service) CountCountry() []model.Item {
	if cache.Exists("count_country") {
		return cache.GetItems("count_country")
	}

	servers := s.store.GetAllServers()

	m := make(map[string]int)
	for _, server := range servers {
		m[server.Country]++
	}
	result := toItems(m)[:min(len(m), 20)]

	cache.SetItems("count_country", result)

	return result
}

// Return the distribution of cooperative/social/madness/competitive servers
// There seems to others (survival, endless, wilderness) on this so only return top 4
func (s Service) CountIntent() []model.Item {
	if cache.Exists("count_intent") {
		return cache.GetItems("count_intent")
	}

	servers := s.store.GetAllServers()

	m := make(map[string]int)
	for _, server := range servers {
		m[server.Intent]++
	}

	result := toItems(m)[:min(len(m), 4)]

	cache.SetItems("count_intent", result)

	return result
}

//
func (s Service) CountPlatform() []model.Item {
	if cache.Exists("count_platform") {
		return cache.GetItems("count_platform")
	}

	servers := s.store.GetAllServers()

	m := make(map[string]int)
	for _, server := range servers {
		m[server.Platform]++
	}

	result := toItems(m)

	cache.SetItems("count_platform", result)

	return result
}

func (s Service) CountSeason() []model.Item {
	if cache.Exists("count_season") {
		return cache.GetItems("count_season")
	}

	servers := s.store.GetAllServers()

	m := make(map[string]int)
	for _, server := range servers {
		m[server.Season]++
	}

	result := toItems(m)[:min(len(m), 4)]

	cache.SetItems("count_season", result)

	return result
}

func (s Service) CountModded() []model.Item {
	if cache.Exists("count_modded") {
		return cache.GetItems("count_modded")
	}

	servers := s.store.GetAllServers()

	m := make(map[string]int)
	for _, server := range servers {
		if server.Mods {
			m["Modded"]++
		} else {
			m["Vanilla"]++
		}
	}

	result := toItems(m)

	cache.SetItems("count_modded", result)

	return result
}

// GetCountries returns a list of all countries of which don't starve servers exist
func (s Service) GetCountries() []string {
	servers := s.store.GetAllServers()

	m := make(map[string]bool)
	for _, server := range servers {
		m[server.Country] = true
	}

	arr := make([]string, 0)
	for key := range m {
		arr = append(arr, key)
	}

	return arr
}
