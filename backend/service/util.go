package service

import (
	"dontstarve-stats/model"
	"sort"
)

// convert map to chart.js data format (array of tuples) format
func toItems(m map[string]int) []model.Item {
	sl := make([]model.Item, 0)

	for key, value := range m {
		sl = append(sl, model.Item{key, value})
	}

	return sortDescending(sl)
}

// convert map to chart.js data format for multi-line line-chart
// [{name: "xyz", data: [Item1, Item2, Item3]}, {name: "abc", data: [Item1, Item2, Item3]}]
func toSeries(m map[string][]model.Item) []model.Series {
	sl := make([]model.Series, 0)

	for name, data := range m {
		sl = append(sl, model.Series{Name: name, Data: data})
	}

	return sl
}

// do not make client sort the items
func sortDescending(slice []model.Item) []model.Item {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i][1].(int) > slice[j][1].(int)
	})
	return slice
}

// go has no generics nice
func sortFloatDescending(slice []model.Item) []model.Item {
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
