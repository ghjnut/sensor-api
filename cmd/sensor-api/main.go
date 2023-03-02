package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	// "context"

	log "github.com/sirupsen/logrus"

	"ghjnut/sensor"
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

	http.HandleFunc("/ingest", ingestHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

type PayloadRaw struct {
	data []string
}

// TODO probably move to a handler struct
func ingestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		bad_payload_handler(w, errors.New("bad request type 'POST'"))
		return
	}

	// TODO this could fail part-way through. we probably should separate parsing and saving
	dec := json.NewDecoder(req.Body)
	var pr PayloadRaw
	err := dec.Decode(&pr)
	if err != nil {
		bad_payload_handler(w, err)
		return
	}

	for i := 0; i < len(pr.data); i++ {
		r, err := sensor.NewReading(pr.data[i])
		if err != nil {
			bad_payload_handler(w, err)
			return
		}

		err = r.Save()
		if err != nil {
			bad_payload_handler(w, err)
			return
		}
	}
}

// what to return if bad decode? should surface, as we would _never_ expect to receive bad payloads
func bad_payload_handler(w http.ResponseWriter, err error) {
	log.Error(err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func init_db() *sql.DB {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
