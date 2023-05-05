package storage

import (
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ResamVi/dontstarve-datavis/alert"
	"github.com/ResamVi/dontstarve-datavis/model"
)

type Store interface {
	GetAge() time.Time
	Start() // Call to keep track of the age of the data

	//CreatePlayer() Not needed because CreateServers creates Players as well
	GetAllPlayers() []model.Player
	DeleteAllPlayers()

	CreateServers(servers []model.Server)
	DeleteAllServers()
	GetAllServers() []model.Server

	GetPlayerStat(name string) model.PlayerStat
	UpdatePlayerStat(model.PlayerStat)

	CreateContinentSnapshot(continents map[string]int)
	GetSeriesContinents() []model.ContinentSnapshot

	CreateCharacterSnapshot(characters map[string]int)
	GetSeriesCharacters() []model.CharacterSnapshot

	CreatePercentageSnapshot(ranking model.PercentageSnapshot)
	GetPercentageSnapshot(character string) []model.PercentageSnapshot
}

type Gorm struct {
	db *gorm.DB
}

func New(url string) Store {
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // Error, Warn, Info
		})

		if err == nil {
			break
		}

		log.Infof("Retrying database connection in 5s (%d/5)", i+1)
		log.Infoln(err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Panicf("Could not connect do database: %s", err.Error())
	}
	log.Info("Connection to database successful")

	db.AutoMigrate(
		&model.Server{},
		&model.Player{},
		&model.ContinentSnapshot{},
		&model.CharacterSnapshot{},
		&model.Character{},
		&model.PercentageSnapshot{},
		&model.PlayerStat{},
		&model.Start{},
	)

	return Gorm{
		db: db,
	}
}

func (gorm Gorm) Start() {
	gorm.db.Create(&model.Start{Started: time.Now()})
}

func (gorm Gorm) GetAge() time.Time {
	var s model.Start
	gorm.db.First(&s)
	return s.Started
}

func (gorm Gorm) GetPlayerStat(name string) model.PlayerStat {
	var stat model.PlayerStat
	gorm.db.FirstOrInit(&stat, model.PlayerStat{Name: name})
	return stat
}

func (gorm Gorm) UpdatePlayerStat(stat model.PlayerStat) {
	gorm.db.Save(&stat)
}

func (gorm Gorm) GetAllPlayers() []model.Player {
	var players []model.Player

	err := gorm.db.Find(&players).Error
	if err != nil {
		panic(err)
	}

	return players
}

func (gorm Gorm) DeleteAllPlayers() {
	gorm.db.Where("1 = 1").Delete(&model.Player{})
}

func (gorm Gorm) CreateServers(servers []model.Server) {
	err := gorm.db.CreateInBatches(servers, 100).Error
	if err != nil {
		alert.String("failed CreateServer")
		panic(err)
	}
}

func (gorm Gorm) GetAllServers() []model.Server {
	var servers []model.Server

	err := gorm.db.Find(&servers).Error
	if err != nil {
		panic(err)
	}

	return servers
}

func (gorm Gorm) DeleteAllServers() {
	gorm.db.Where("1 = 1").Delete(&model.Server{})
}

func (gorm Gorm) CreateContinentSnapshot(continents map[string]int) {
	gorm.db.Create(&model.ContinentSnapshot{
		Date:         time.Now().Local(),
		Asia:         continents["Asia"],
		Europe:       continents["Europe"],
		NorthAmerica: continents["North America"],
		SouthAmerica: continents["South America"],
		Africa:       continents["Africa"],
		Oceania:      continents["Oceania"],
	})
}

func (gorm Gorm) CreateCharacterSnapshot(characters map[string]int) {
	chars := make([]model.Character, 0)
	for name, count := range characters {
		chars = append(chars, model.Character{Name: name, Count: count})
	}

	gorm.db.Create(&model.CharacterSnapshot{
		Date:       time.Now().Local(),
		Characters: chars,
	})
}

func (gorm Gorm) CreatePercentageSnapshot(ranking model.PercentageSnapshot) {
	gorm.db.Create(&ranking)
}

func (gorm Gorm) GetSeriesContinents() []model.ContinentSnapshot {
	var continents []model.ContinentSnapshot
	gorm.db.Find(&continents)
	return continents
}

func (gorm Gorm) GetSeriesCharacters() []model.CharacterSnapshot {
	var characters []model.CharacterSnapshot
	gorm.db.Preload("Characters").Find(&characters)
	return characters
}

func (gorm Gorm) GetPercentageSnapshot(name string) []model.PercentageSnapshot {
	var rankings []model.PercentageSnapshot
	gorm.db.Where("character = ?", name).Find(&rankings)
	return rankings
}
