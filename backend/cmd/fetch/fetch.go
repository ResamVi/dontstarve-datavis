package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/ResamVi/dontstarve-datavis/fetch/fetcher"
	"github.com/ResamVi/dontstarve-datavis/model"

	"github.com/olekukonko/tablewriter"
)

func main() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Africa", "Empty", "Europe", "ERROR", "South America", "Asia", "North America", "Oceania"})
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	go func() {
		<-sc
		table.Render()
		os.Exit(0)
	}()

	for {

		// B Processing
		// data, err := fetcher.GetServers()
		// if err != nil {
		// 	panic(err)
		// }
		// for _, entry := range data {
		// 	table.Append([]string{entry.Addr, entry.Name})
		// }
		// table.Render()

		// A Processing
		// table.SetHeader([]string{"Name", "Sign", "Rating"}) // A
		f := fetcher.New()
		data, err := f.Fetch()
		if err != nil {
			panic(err)
		}
		count := make(map[string]int)

		for _, entry := range data {
			for _, player := range entry.Players {
				count[player.Continent]++
			}
		}

		// C Processing
		table.Append(conv([]int{count["Africa"], count[""], count["ERROR"], count["South America"], count["Asia"], count["North America"], count["Oceania"]}))

		// for k, v := range count {
		// 	fmt.Printf("%v: %v\n", k, v)
		// }

		fmt.Println("----------------------------------")
	}
}

func playerStr(players []model.Player) string {
	var str []string

	for _, player := range players {
		line := fmt.Sprintf("%v (%v)", player.Name, player.Country)
		str = append(str, line)
	}

	return strings.Join(str, ", ")
}

func conv(nums []int) []string {
	result := make([]string, len(nums))
	for i := range nums {
		result[i] = strconv.Itoa(nums[i])
	}
	return result
}
