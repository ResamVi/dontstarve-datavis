package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	host     = "localhost" //os.Getenv("HOST")
	port     = 5432
	user     = "root"       //os.Getenv("POSTGRES_USER")
	password = "password"   //os.Getenv("POSTGRES_PASSWORD")
	dbname   = "mydatabase" //os.Getenv("POSTGRES_DB")
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	creds := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", creds)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/attribute/:attr", func(c *gin.Context) {

		var rows *sql.Rows

		switch c.Param("attr") {
		case "intent":
			rows, err = db.Query("SELECT * FROM count_intent LIMIT 4")
		case "platforms":
			rows, err = db.Query("SELECT * FROM count_platform")
		case "modded":
			rows, err = db.Query("SELECT * FROM count_vanilla")
		case "season":
			rows, err = db.Query("SELECT * FROM count_season")
		}

		handleError(err)
		list := make(map[string]int)
		for rows.Next() {
			var key string
			var value int

			err = rows.Scan(&key, &value)
			handleError(err)

			list[key] = value
		}

		c.JSON(200, list)
	})

	router.GET("/count/:type", func(c *gin.Context) {

		var row *sql.Row

		switch c.Param("type") {
		case "players":
			row = db.QueryRow("SELECT player_count FROM count")
		case "servers":
			row = db.QueryRow("SELECT server_count FROM count")
		default:
			panic("Wrong type")
		}

		var count int
		err := row.Scan(&count)
		handleError(err)

		c.JSON(200, count)
	})

	router.GET("/countries", func(c *gin.Context) {

		var rows *sql.Rows
		rows, err = db.Query("SELECT origin FROM count_player")

		handleError(err)

		list := make([]string, 5)
		for rows.Next() {
			var countryName string

			err = rows.Scan(&countryName)
			handleError(err)

			list = append(list, countryName)
		}

		c.JSON(200, list)
	})

	router.GET("/origin/:type", func(c *gin.Context) {

		var rows *sql.Rows

		switch c.Param("type") {
		case "players":
			rows, err = db.Query("SELECT * FROM count_player LIMIT 20")
		case "servers":
			rows, err = db.Query("SELECT * FROM count_server LIMIT 20")
		default:
			panic("Wrong type")
		}

		handleError(err)

		list := make(map[string]int)
		for rows.Next() {
			var countryName string
			var count int

			err = rows.Scan(&countryName, &count)
			handleError(err)

			list[countryName] = count
		}

		c.JSON(200, list)
	})

	router.GET("/characters/:origin", func(c *gin.Context) {
		// TODO: Sanitize
		rows, err := db.Query("SELECT character, count FROM count_character_by_origin WHERE origin = $1 LIMIT 20", strings.Title(c.Param("origin")))
		handleError(err)

		characters := make(map[string]int)
		for rows.Next() {
			var characterName string
			var count int

			err = rows.Scan(&characterName, &count)
			handleError(err)

			switch characterName {
			case "":
				characterName = "pending"
			case "wathgrithr":
				characterName = "wigfrid"
			case "waxwell":
				characterName = "maxwell"
			case "monkey_king":
				characterName = "wilbur"
			}

			characters[characterName] = count
		}

		c.JSON(200, characters)
	})

	router.GET("/characters", func(c *gin.Context) {

		rows, err := db.Query("SELECT * FROM count_character ORDER BY count DESC LIMIT 19")
		handleError(err)

		characters := make(map[string]int)
		for rows.Next() {
			var characterName string
			var count int

			err = rows.Scan(&characterName, &count)
			handleError(err)

			switch characterName {
			case "":
				characterName = "pending"
			case "wathgrithr":
				characterName = "wigfrid"
			case "waxwell":
				characterName = "maxwell"
			case "monkey_king":
				characterName = "wilbur"
			}

			characters[characterName] = count
		}

		c.JSON(200, characters)
	})
	router.Run(":3000")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
