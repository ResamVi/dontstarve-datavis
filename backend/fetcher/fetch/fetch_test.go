package fetch

import (
	"os"
	"testing"
)

func initFetch() (*Fetch, error) {
	fetch, err := Init("")
	return fetch, err
}

func TestToken(t *testing.T) {
	os.Setenv("TOKEN", "YOUR TOKEN HERE")

	f, err := initFetch()
	if err != nil {
		t.Errorf(err.Error())
	}

	servers, err := f.readServerList()
	if len(servers) == 0 || err != nil {
		t.Errorf("Token test failed: " + err.Error())
	}
}

func TestGeolocate(t *testing.T) {
	f, err := initFetch()
	if err != nil {
		t.Errorf(err.Error())
	}

	country, continent, iso := f.geolocate("91.132.146.174")
	if country != "Germany" || continent != "Europe" || iso != "DE" {
		t.Error("Wrong answer")
	}

	country, continent, iso = f.geolocate("127.0.0.1")
	if country != "" || continent != "" || iso != "" {
		t.Error("Wrong answer")
	}
}
