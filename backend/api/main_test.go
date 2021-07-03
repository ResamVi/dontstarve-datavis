package main

import (
	"dontstarve-stats/fetch/fetcher"
	"dontstarve-stats/service"
	"dontstarve-stats/storage"
	"testing"
	"time"
)

func setup() service.Service {
	store := storage.New(
		"localhost",
		"root",
		"password",
		"mydatabase",
		"5432",
	)

	fetch := fetcher.New("<TOKEN HERE>")
	servers := fetch.Servers()

	svc := service.New(store)
	svc.TrackPlaytime(servers[:1000], time.Now())
	svc.ServerSnapshot(servers)
	svc.ContinentSnapshot()
	svc.CharacterSnapshot()
	svc.PercentageSnapshot()

	return svc
}

// 270124390 ns/op --> 308616 ns/op
func BenchmarkCount(b *testing.B) {
	svc := setup()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		svc.CountServers()
		svc.CountPlayers()
	}
}

// 571416866 ns/op	51720980 B/op	 1558939 allocs/op
// --> 332610 ns/op	    7784 B/op	     200 allocs/op
func BenchmarkCharacterCountry(b *testing.B) {
	svc := setup()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		svc.CountCharacters(false)
		svc.CountCountry()
	}
}

func BenchmarkOriginCountries(b *testing.B) {
	svc := setup()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		svc.CountPlayerOrigin(false)
		svc.GetCountries()
	}
}
