package klei

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Regioncapabilities is Klei's HTTP response from /regioncapabilities-v2.json
//
// Example:
//
//	{
//	    "LobbyRegions": [
//	        { "Region": "us-east-1" },
//	        { "Region": "eu-central-1" },
//	        { "Region": "ap-southeast-1" },
//	        { "Region": "ap-east-1" }
//	    ]
//	}
type Regioncapabilities struct {
	LobbyRegions []struct {
		Region string `json:"Region"`
	} `json:"LobbyRegions"`
}

// Thanks to Crestwave for noting what changed.
// https://forums.kleientertainment.com/forums/topic/138537-march-quality-of-life-update-now-live/?do=findComment&comment=1551936
const regioncapabilitiesURL = "https://lobby-v2-cdn.klei.com/regioncapabilities-v2.json"

// Regions lists all AWS regions that Klei hosts their servers on.
func Regions() ([]string, error) {
	var response Regioncapabilities

	resp, err := http.Get(regioncapabilitiesURL)
	if err != nil {
		return nil, fmt.Errorf("could not send GET request to %v: %w", regioncapabilitiesURL, err)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("could not decode GET response body from %v: %w", regioncapabilitiesURL, err)
	}

	// Convert to []string for easier usage.
	var result []string
	for _, regions := range response.LobbyRegions {
		result = append(result, regions.Region)
	}

	return result, nil

}
