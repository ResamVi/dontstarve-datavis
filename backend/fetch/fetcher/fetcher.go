package fetcher

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"

	"github.com/ResamVi/dontstarve-datavis/fetch/klei"
	"github.com/ResamVi/dontstarve-datavis/model"
)

// Fetcher is the main orchestrator to get the data
type Fetcher struct {
	cycle int // how often Fetch() has been called
	geoip *geoip2.Reader
}

//go:embed GeoLite2-Country.mmdb
var geoipfile []byte

// New runs all necessary steps to request, retrieve, interpret and persist don't starve together server info
func New() *Fetcher {
	geoip, err := geoip2.FromBytes(geoipfile)
	if err != nil {
		panic("could not open geolite db: " + err.Error())
	}

	return &Fetcher{geoip: geoip}
}

// Fetch reads and parses all servers (and their containing players) shown in the server list
//
// If it cannot do any of those things it will fail silently: this
// behavior is intended to keep running but simply skip one cycle.
func (f *Fetcher) Fetch() ([]model.Server, error) {
	servers, err := getServers()
	if err != nil {
		return nil, err
	}

	result := make([]model.Server, 0)
	for _, entry := range servers {
		country, continent, iso := entry.Location(f.geoip)

		server := model.Server{
			Country:   country,
			Iso:       iso,
			Continent: continent,

			Platform: klei.ParsePlatform(entry.Platform),
			// Elapsed:   f.parseElapsed(server["data"].(string)),
			Cycle: f.cycle,

			Name:      entry.Name,
			Connected: entry.Connected,
			Mode:      entry.Mode,
			Season:    entry.Season,
			Intent:    entry.Intent,
			Mods:      entry.Mods,
			Date:      time.Now().Local(),
			Players:   klei.ParsePlayers(entry.Players, country, continent, iso),
		}

		result = append(result, server)
	}

	f.cycle++

	return result, nil
}

// getServers reads from all of klei's endpoints.
func getServers() ([]klei.Server, error) {
	regions, err := klei.Regions()
	if err != nil {
		return nil, fmt.Errorf("could not get regions: %w", err)
	}

	var result []klei.Server
	for _, region := range regions {
		log.Infof("Parsing %v", region)

		start := time.Now()

		// Retrieve the list of servers visible in the lobby.
		lobby, err := klei.Lobby(region)
		if err != nil {
			continue
		}

		log.Infof("Found %v servers", len(lobby))

		// Collect the details (player info etc.) of each server.
		servers, count := klei.Servers(lobby, region)
		result = append(result, servers...)

		log.Infof("Parsed in %.2fs with %v/%v failing", time.Since(start).Seconds(), count, len(servers))
	}

	return result, nil
}

// first number indicates days elapsed
var regElapsed = regexp.MustCompile(`\d+`)

// parseElapsed converts the day elapsed counter to int
func (f *Fetcher) parseElapsed(data string) int {
	str := regElapsed.FindString(data)
	if str == "" { // TODO: Found anomaly
		log.Debug("Found server without proper data field")
		return 0
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		panic("Could not convert: '" + str + "'")
	}

	return i
}

// Cycle returns how often we fetched from klei
func (f *Fetcher) Cycle() int {
	return f.cycle
}
