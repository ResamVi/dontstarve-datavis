package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"dontstarve-stats/alert"
	"dontstarve-stats/fetcher/fetch"
	"dontstarve-stats/service"

	"dontstarve-stats/storage"
)

func main() {
	defer alert.Msg("Service 'fetcher' has stopped running")

	var store storage.Store
	if !isProd() {
		store = storage.Init(
			"localhost",
			"root",
			"password",
			"mydatabase",
			"5432",
		)

		os.Setenv("TOKEN", "<INSERT YOUR TOKEN>")
	} else {
		store = storage.Init(
			"db",
			os.Getenv("USER"),
			os.Getenv("PASSWORD"),
			os.Getenv("DBNAME"),
			os.Getenv("DBPORT"),
		)
	}

	fetch, err := fetch.Init(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	svc := service.New(store)

	log.Info("Start fetching from Klei's servers...")

	previousCycle := time.Now()
	for {
		start := time.Now()
		servers := fetch.Servers()
		log.Infof("%d servers fetched in %.2fs", len(servers), time.Since(start).Seconds())

		svc.TrackPlaytime(servers, previousCycle)

		start = time.Now()
		svc.ServerSnapshot(servers)
		svc.ContinentSnapshot()
		svc.CharacterSnapshot()
		svc.PercentageSnapshot()
		log.Infof("Finished persisting in %.2fs", time.Since(start).Seconds())

		log.Infof("Finished cycle no. %d", fetch.Cycle())
		previousCycle = time.Now()

		time.Sleep(15 * time.Minute)
	}
}

// docker run --rm --name postgrestmp --network="host" -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e POSTGRES_DB=mydatabase postgres

func isProd() bool {
	_, isProd := os.LookupEnv("PROD")
	return isProd
}
