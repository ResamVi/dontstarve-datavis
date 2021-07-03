package fetcher

import (
	"os"
	"testing"
)

func TestToken(t *testing.T) {
	os.Setenv("TOKEN", "YOUR TOKEN HERE")

	f := New("")

	servers, err := f.readServerList()
	if len(servers) == 0 || err != nil {
		t.Errorf("Token test failed: " + err.Error())
	}
}

func TestGeolocate(t *testing.T) {
	f := New("")

	country, continent, iso := f.geolocate("91.132.146.174")
	if country != "Germany" || continent != "Europe" || iso != "DE" {
		t.Error("Wrong answer")
	}

	country, continent, iso = f.geolocate("127.0.0.1")
	if country != "" || continent != "" || iso != "" {
		t.Error("Wrong answer")
	}
}
