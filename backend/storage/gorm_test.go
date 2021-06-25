package storage

import (
	"os"
	"testing"
)

func initDB() Store {
	os.Setenv("HOST", "localhost")
	os.Setenv("USER", "root")
	os.Setenv("PASSWORD", "password")
	os.Setenv("DBNAME", "mydatabase")
	os.Setenv("DBPORT", "5432")
	os.Setenv("TOKEN", "<YOUR TOKEN HERE>")

	return Init(
		os.Getenv("HOST"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPORT"),
	)
}

func TestDeleteAllPlayers(t *testing.T) {
	//store := initDB()

	/*if err != nil {
		t.Errorf(err)
	}*/
}
