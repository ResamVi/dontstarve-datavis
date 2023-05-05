package klei

import (
	"fmt"

	"github.com/ResamVi/dontstarve-datavis/alert"
)

// ID conversion table for platforms (they call them "platforms")
// (see: https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data/?do=findComment&comment=1306033)
var platformID = map[int]string{
	1:  "Steam",
	2:  "PSN",
	4:  "WeGame",
	10: "XBOX LIVE",
	16: "16",
	19: "19",
	32: "32",
}

// platforms (they call them "platforms") on which Don't Starve Together is played:
//
// Taken from:
// https://forums.kleientertainment.com/forums/topic/145450-new-server-lobby-and-networking-beta-536845/?do=findComment&comment=1616725
var platforms = []string{
	"Steam",
	"PSN",
	"Rail",
	"XBone",
	"Switch",
}

// ParsePlatform convert platform number to its name.
func ParsePlatform(platformNumber int) string {
	fmt.Println(platformNumber)
	if _, exists := platformID[platformNumber]; !exists {
		alert.Stringf("unknown: %v", platformNumber)
	}
	return platformID[platformNumber]
}
