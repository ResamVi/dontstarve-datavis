package fetcher

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"dontstarve-stats/alert"
	"dontstarve-stats/model"

	"github.com/oschwald/geoip2-golang"
)

// Fetcher is the main orchestrator to get the data
type Fetcher struct {
	cycle int    // how often Servers() has been called
	token string // token as per https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data
	geoip *geoip2.Reader
}

// ParsedResponse is Klei server's answer to our post request.
// Contains single key "GET" that points to the server list
type ParsedResponse map[string]ServerList

// ServerList contains *all* available data of servers that are registered to this endpoint
type ServerList []ServerJSON

// ServerJSON is json representation of a single joinable Don't Starve Together ServerJSON
type ServerJSON map[string]interface{}

//go:embed GeoLite2-Country.mmdb
var geoipfile []byte

// New runs all necessary steps to request, retrieve, interpret and persist don't starve together server info
func New(token string) *Fetcher {
	geoip, err := geoip2.FromBytes(geoipfile)
	if err != nil {
		panic("could not open geolite db: " + err.Error())
	}

	return &Fetcher{token: token, geoip: geoip}
}

// Servers reads and parses all servers (and their containing players) shown in the server list
//
// If it cannot do any of those things it will fail silently: this
// behavior is intended to keep running but simply skip one cycle
func (f *Fetcher) Servers() []model.Server {
	serverList, err := f.readServerList()
	if err != nil {
		alert.Msg(err.Error())
	}

	servers := make([]model.Server, 0)
	for _, entry := range serverList {
		server := f.parseServer(entry)
		players := f.parsePlayers(entry)

		server.Players = players
		servers = append(servers, server)
	}

	f.cycle++

	return servers
}

// parseServer parses the data and reduces it down the relevant data we require
func (f Fetcher) parseServer(server ServerJSON) model.Server {
	country, continent, iso := f.geolocate(server["__addr"].(string))

	return model.Server{
		Country:   country,
		Iso:       iso,
		Continent: continent,
		Platform:  f.parsePlatform(server["platform"].(float64)),
		Elapsed:   f.parseElapsed(server["data"].(string)),
		Cycle:     f.cycle,
		Name:      server["name"].(string),
		Connected: server["connected"].(float64),
		Mode:      server["mode"].(string),
		Season:    server["season"].(string),
		Intent:    server["intent"].(string),
		Mods:      server["mods"].(bool),
		Date:      time.Now().Local(),
	}
}

// parsePlayer creates the player that are part of a server
func (f Fetcher) parsePlayers(server ServerJSON) []model.Player {
	country, continent, iso := f.geolocate(server["__addr"].(string)) // TODO: redundant: done twice

	if f.isEmpty(server) {
		return []model.Player{}
	}

	r := strings.NewReplacer(
		`return {`, "[",
		`colour=`, `"colour": `,
		`["colour"]=`, `"colour": `,
		`eventlevel=`, `"eventlevel": `,
		`["eventlevel"]=`, `"eventlevel": `,
		`name=`, `"name": `,
		`["name"]=`, `"name": `,
		`netid=`, `"netid": `,
		`["netid"]=`, `"netid": `,
		`prefab=`, `"character": `,
		`["prefab"]=`, `"character": `,
		"\n", "",
		"\t", "",
		"\a", "",
		"\x11", "",
		"\x01", "",
		"\x14", "",
		"\x0e", "",
		"ó", "o",
	)
	s := r.Replace(server["players"].(string))

	b := []byte(s)
	b[len(b)-1] = ']'

	var ps []model.Player
	err := json.Unmarshal(b, &ps)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			alert.Msg(fmt.Sprintf("unmarshal error at byte offset %d", e.Offset))
		}

		return []model.Player{}
	}

	for i := range ps {
		ps[i].Character = f.translate(ps[i].Character)
		ps[i].Continent = continent
		ps[i].Country = country
		ps[i].Iso = iso
	}

	return ps
}

// geolocate uses GeoIP to get a country origin of the IP
func (f Fetcher) geolocate(ip string) (string, string, string) {
	record, err := f.geoip.Country(net.ParseIP(ip))
	if err != nil {
		return "ERROR", "ERROR", "ERROR" // TODO: What to do?
	}

	return record.Country.Names["en"], record.Continent.Names["en"], record.Country.IsoCode
}

// some charaacters are nicknamed differently
// also capitalize first letter for easy display
func (f Fetcher) translate(name string) string {
	switch name {
	case "":
		return "<Selecting>"
	case "wathgrithr":
		return "Wigfrid"
	case "waxwell":
		return "Maxwell"
	case "wx78":
		return "WX-78"
	case "monkey_king":
		return "Wilbur"
	}

	return strings.Title(name)
}

func (f Fetcher) isEmpty(server ServerJSON) bool {
	return server["players"] == "return {  }"
}

// Klei's endpoint URLs to get the data
var endpoints = []string{
	"https://lobby-us.kleientertainment.com/lobby/read",
	"https://lobby-eu.kleientertainment.com/lobby/read",
	"https://lobby-china.kleientertainment.com/lobby/read",
	"https://lobby-sing.kleientertainment.com/lobby/read",
}

// readServerList reads from all of klei's endpoints
func (f Fetcher) readServerList() (ServerList, error) {
	payload := fmt.Sprintf(`{
		"__token": "%s", 
		"__gameId": "DST", 
		"query": {}
	}`, f.token)

	serverlist := make([]ServerJSON, 0)
	for _, endpoint := range endpoints {

		resp, err := http.Post(endpoint, "application/json", strings.NewReader(payload))
		if err != nil {
			return []ServerJSON{}, errors.New("could not request server list: " + err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []ServerJSON{}, errors.New("Could not read answer to request: " + err.Error() + "\n contents: " + string(body))
		}

		if string(body) == `{"error":"AUTH_ERROR_E_EXPIRED_TOKEN"}` {
			return []ServerJSON{}, errors.New("authorization failed, token seems to be expired/invalid")
		}

		var servers ParsedResponse
		err = json.Unmarshal(body, &servers)
		if err != nil {
			if e, ok := err.(*json.SyntaxError); ok {
				alert.Msg(fmt.Sprintf("syntax error at byte offset %d", e.Offset))
			}

			return []ServerJSON{}, errors.New("could not unmarshal answer: '" + string(body) + "'\n" + err.Error())
		}

		serverlist = append(serverlist, servers["GET"]...)
	}

	return serverlist, nil
}

// first number indicates days elapsed
var regElapsed = regexp.MustCompile(`\d+`)

// parseElapsed converts the day elapsed counter to int
func (f Fetcher) parseElapsed(data string) int {
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

// Conversion table for platforms
// (see: https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data/?do=findComment&comment=1306033)
var platform = map[float64]string{
	1:  "Steam",
	2:  "PSN",
	4:  "WeGame",
	10: "XBOX LIVE",
	16: "???",
	19: "???",
	32: "???",
}

// Convert platform number to name
func (f Fetcher) parsePlatform(platformNumber float64) string {
	if _, exists := platform[platformNumber]; !exists {
		log.Panicf("unknown: %f", platformNumber)
	}
	return platform[platformNumber]
}

// Cycle returns how often we fetched from klei
func (f Fetcher) Cycle() int {
	return f.cycle
}
