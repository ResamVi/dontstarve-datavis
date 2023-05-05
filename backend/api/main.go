package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ResamVi/dontstarve-datavis/alert"
	"github.com/ResamVi/dontstarve-datavis/service"
	"github.com/ResamVi/dontstarve-datavis/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// go get -u github.com/gin-gonic/gin
func main() {
	defer alert.String("Service 'api' has stopped running")

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

	svc := service.New(storage.New(url))

	// CORS
	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	r := gin.Default()
	r.Use(cors.New(config))

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "0.1")
	})

	r.GET("/meta/servers", func(c *gin.Context) {
		c.JSON(200, svc.CountServers())
	})

	r.GET("/meta/players", func(c *gin.Context) {
		c.JSON(200, svc.CountPlayers())
	})

	r.GET("/characters", func(c *gin.Context) {
		includeModded := c.DefaultQuery("modded", "false") == "true"
		c.JSON(200, svc.CountCharacters(includeModded))
	})

	r.GET("/count/country", func(c *gin.Context) {
		c.JSON(200, svc.CountCountry())
	})

	r.GET("/count/players", func(c *gin.Context) {
		includeAll := c.DefaultQuery("all", "false") == "true"
		c.JSON(200, svc.CountPlayerOrigin(includeAll))
	})

	r.GET("/meta/countries", func(c *gin.Context) {
		c.JSON(200, svc.GetCountries())
	})

	r.GET("/meta/age", func(c *gin.Context) {
		c.JSON(200, svc.LastUpdate())
	})

	r.GET("/characters/:country", func(c *gin.Context) {
		name := c.Param("country")
		c.JSON(200, svc.GetCountryPreference(name))
	})

	r.GET("/series/continents", func(c *gin.Context) {
		c.JSON(200, svc.GetSeriesContinents())
	})

	r.GET("/series/characters", func(c *gin.Context) {
		c.JSON(200, svc.GetSeriesCharacters())
	})

	r.GET("/characters/percentage/:character", func(c *gin.Context) {
		character := c.Param("character")
		c.JSON(200, svc.GetCountryCharacters(character))
	})

	r.GET("/count/intent", func(c *gin.Context) {
		c.JSON(200, svc.CountIntent())
	})

	r.GET("/count/platforms", func(c *gin.Context) {
		c.JSON(200, svc.CountPlatform())
	})

	r.GET("/count/season", func(c *gin.Context) {
		c.JSON(200, svc.CountSeason())
	})

	r.GET("/count/modded", func(c *gin.Context) {
		c.JSON(200, svc.CountModded())
	})

	r.GET("/series/preferences/:character", func(c *gin.Context) {
		name := c.Param("character")
		c.JSON(200, svc.GetCountryRankings(name))
	})

	r.GET("/player/character/:player", func(c *gin.Context) {
		name := c.Param("player")
		c.JSON(200, svc.GetPlayTime(name))
	})

	r.GET("/meta/started", func(c *gin.Context) {
		c.JSON(200, svc.Started())
	})

	r.Run("0.0.0.0:8003")
}

// isProd assumes we run in a docker environment
func isProd() bool {
	_, isProd := os.LookupEnv("PROD")
	return isProd
}
