package fetcher

import (
	"os"
	"testing"
)

func TestToken(t *testing.T) {
	os.Setenv("TOKEN", "YOUR TOKEN HERE")

	servers, err := GetServers()
	if len(servers) == 0 || err != nil {
		t.Errorf("Token test failed: " + err.Error())
	}
}
