package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ResamVi/dontstarve-datavis/alert"
	"github.com/ResamVi/dontstarve-datavis/fetch/fetcher"
	"github.com/ResamVi/dontstarve-datavis/service"
	"github.com/ResamVi/dontstarve-datavis/storage"
)

func main() {
	defer alert.String("Service 'fetcher' has stopped running")

	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASSWORD"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
	)

	// Default if nothing is set.
	if _, exists := os.LookupEnv("DBUSER"); !exists {
		url = "postgres://root:password@localhost:5432/dststats"
	}

	fetch := fetcher.New()

	svc := service.New(storage.New(url))

	log.Info("Start fetching from Klei's servers...")

	previousCycle := time.Now()
	for {
		start := time.Now()
		servers, err := fetch.Fetch()
		if err != nil {
			alert.Error(err)
			continue
		}

		if len(servers) == 0 {
			alert.String("No servers found")
			continue
		}

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
