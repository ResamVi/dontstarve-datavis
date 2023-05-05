package klei

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gammazero/workerpool"
	"github.com/oschwald/geoip2-golang"
)

type Server struct {
	// These fields are always filled.
	Addr            string                 `json:"__addr"`
	RowID           string                 `json:"__rowId"`
	Host            string                 `json:"host"`
	Clanonly        bool                   `json:"clanonly"`
	Platform        int                    `json:"platform"`
	Mods            bool                   `json:"mods"`
	Name            string                 `json:"name"`
	Pvp             bool                   `json:"pvp"`
	Session         string                 `json:"session"`
	Fo              bool                   `json:"fo"`
	Password        bool                   `json:"password"`
	Guid            string                 `json:"guid"`
	Maxconnections  int                    `json:"maxconnections"`
	Dedicated       bool                   `json:"dedicated"`
	Clienthosted    bool                   `json:"clienthosted"`
	Connected       float64                `json:"connected"`
	Mode            string                 `json:"mode"`
	Port            int                    `json:"port"`
	V               int                    `json:"v"`
	Tags            string                 `json:"tags"`
	Season          string                 `json:"season"`
	Lanonly         bool                   `json:"lanonly"`
	Intent          string                 `json:"intent"`
	Allownewplayers bool                   `json:"allownewplayers"`
	Serverpaused    bool                   `json:"serverpaused"`
	Steamid         string                 `json:"steamid"`
	Steamroom       string                 `json:"steamroom"`
	Secondaries     map[string]interface{} `json:"secondaries"`

	Data          string        `json:"data"`
	Worldgen      string        `json:"worldgen"`
	Players       string        `json:"players"`
	ModsInfo      []interface{} `json:"mods_info"`
	Tick          int           `json:"tick"`
	Clientmodsoff bool          `json:"clientmodsoff"`
	Nat           int           `json:"nat"`
}

// Location infers more geological details of the server based on its IP.
func (d Server) Location(geoip *geoip2.Reader) (string, string, string) {
	record, err := geoip.Country(net.ParseIP(d.Addr))
	if err != nil {
		return "ERROR", "ERROR", "ERROR"
	}

	return record.Country.Names["en"], record.Continent.Names["en"], record.Country.IsoCode
}

func Servers(lobby []LobbyEntry, region string) ([]Server, int) {
	var result []Server

	wp := workerpool.New(1000)

	ch := make(chan struct{}, len(lobby))
	for _, entry := range lobby {
		wp.Submit(func() {
			details, err := readEntry(entry, region)
			if err != nil {
				ch <- struct{}{}
				return
			}

			result = append(result, *details)
		})
	}

	wp.StopWait()

	// Atomic counter did not work. Poor try at counting errors.
	count := 0
	for len(ch) > 0 {
		count++
	}

	return result, count
}

// token as per https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data
var template = fmt.Sprintf(`{
	"__gameId": "DST",
	"__token": "%v",
	"query": {
		"__rowId": "%%v"
	}
}`, os.Getenv("TOKEN"))

// readEntry takes an entry and retrieves all server infos.
func readEntry(entry LobbyEntry, region string) (*Server, error) {
	query := fmt.Sprintf(template, entry.RowID)
	url := fmt.Sprintf("https://lobby-v2-%v.klei.com/lobby/read", region)

	resp, err := http.Post(url, "application/json", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("could not send POST request to %v: %w", url, err)
	}

	// The actual data is wrapped in a "GET" json field.
	var response WrappedResponse[[]Server]

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("could not decode POST response body from %v: %w", url, err)
	}

	if len(response.GET) != 1 {
		return nil, fmt.Errorf("unexpected amount of servers received for %v: %v", url, response.GET)
	}

	return &response.GET[0], nil
}
