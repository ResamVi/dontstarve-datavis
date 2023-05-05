package klei

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ResamVi/dontstarve-datavis/alert"
	"github.com/ResamVi/dontstarve-datavis/model"
)

// ParsePlayers creates the player that are part of a server
func ParsePlayers(raw, country, continent, iso string) []model.Player {
	// Server is empty
	if raw == "return {  }" {
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
		"\x10", "",
		"ó", "o",
	)
	s := r.Replace(raw)

	b := []byte(s)
	b[len(b)-1] = ']'

	var ps []model.Player
	err := json.Unmarshal(b, &ps)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			alert.String(e.Error() + " " + string(b))
			alert.String(fmt.Sprintf("unmarshal error at byte offset %d", e.Offset))
		}

		return []model.Player{}
	}

	for i := range ps {
		ps[i].Character = translate(ps[i].Character)
		ps[i].Continent = continent
		ps[i].Country = country
		ps[i].Iso = iso
	}

	return ps
}

// some charaacters are nicknamed differently
// also capitalize first letter for easy display
func translate(name string) string {
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
