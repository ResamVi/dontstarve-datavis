package main

import (
	"database/sql"
	"fmt"

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

	router.GET("/test", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"value": 123,
			"odor":  150,
		},
		)
	})
	router.GET("/characters", func(c *gin.Context) {

		rows, err := db.Query("SELECT * FROM count_character ORDER BY count DESC LIMIT 19")
		handleError(err)

		characters := make(map[string]int)
		for rows.Next() {
			var (
				characterName string
				count         int
			)

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
			fmt.Println(characterName)
		}
		fmt.Println(characters)
		c.JSON(200, characters)
	})
	router.Run(":3000")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
