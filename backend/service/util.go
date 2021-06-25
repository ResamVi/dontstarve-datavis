package service

import "sort"

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

// convert map to chart.js data format (array of tuples) format
func toItems(m map[string]int) []Item {
	sl := make([]Item, 0)

	for key, value := range m {
		sl = append(sl, Item{key, value})
	}

	return sortDescending(sl)
}

// convert map to chart.js data format for multi-line line-chart
// [{name: "xyz", data: [Item1, Item2, Item3]}, {name: "abc", data: [Item1, Item2, Item3]}]
func toSeries(m map[string][]Item) []Series {
	sl := make([]Series, 0)

	for name, data := range m {
		sl = append(sl, Series{Name: name, Data: data})
	}

	return sl
}

// do not make client sort the items
func sortDescending(slice []Item) []Item {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i][1].(int) > slice[j][1].(int)
	})
	return slice
}

// go has no generics nice
func sortFloatDescending(slice []Item) []Item {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i][1].(float64) > slice[j][1].(float64)
	})
	return slice
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
