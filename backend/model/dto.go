package model

// Item is a Tuple.
// Used as a key-value pair that is suited for chart.js to display
// e.g. ["Wigfrid", 15] or ["Russia", 1243]
type Item [2]interface{}

// IsoItem is a Triple
// Used only for "Highest Character Preferences by Country" atm
// e.g. ["Greece", 42.5], ["Poland", 25.5]
type IsoItem [3]interface{}

// Series is a single line of a line-chart
// with `data` having all the dots
type Series struct {
	Name string `json:"name"`
	Data []Item `json:"data"`
}
