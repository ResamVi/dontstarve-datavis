package klei

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LobbyEntry is a single entry in the Lobby that contains data
// used for a followup query to retrieve more server data.
type LobbyEntry struct {
	Addr  string `json:"__addr"`
	RowID string `json:"__rowId"`
}

// Lobby retrieves the listed servers.
func Lobby(region string) ([]LobbyEntry, error) {
	var servers []LobbyEntry

	for _, platform := range platforms {
		url := fmt.Sprintf("https://lobby-v2-cdn.klei.com/%v-%v.json.gz", region, platform)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("could not send GET request to %v: %w", url, err)
		}

		// The actual data is wrapped in a "GET" json field.
		var response WrappedResponse[[]LobbyEntry]

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, fmt.Errorf("could not decode GET response body from %v: %w", url, err)
		}

		servers = append(servers, response.GET...)
	}

	return servers, nil
}
