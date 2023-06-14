package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	// a little dicey with the 'log' namespace collision
	log "github.com/sirupsen/logrus"

	"ghjnut/sensor/internal/httptransport"
	"ghjnut/sensor/internal/service"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func main() {
	// parse config
	debug := true

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO signal handeling goes here

	s := service.New(db)
	h := httptransport.Handler(s)
	err = http.ListenAndServe(":8000", h)
	log.Fatalln(err)
}

func initDB() (*sql.DB, error) {
	host, user, password, dbname := "database", "pguser", "pgpassword", "code_challenge"

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	return sql.Open("postgres", psqlInfo)
}
