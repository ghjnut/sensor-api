package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	// "context"

	_ "github.com/lib/pq"
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

	db, err := init_db()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := &Service{db: db}

	log.Fatal(http.ListenAndServe(":8000", s))
}

type PayloadRaw struct {
	data []string
}

type Service struct {
	db *sql.DB
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case "/ingest":
		s.readingHandler(w, r)
	}
}

func (s *Service) readingHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		bad_payload_handler(w, errors.New("bad request type 'POST'"))
		return
	}

	dec := json.NewDecoder(req.Body)
	var pr PayloadRaw
	err := dec.Decode(&pr)
	if err != nil {
		bad_payload_handler(w, err)
		return
	}

	// TODO this could fail part-way through. we probably should separate parsing and saving
	for i := 0; i < len(pr.Data); i++ {
		log.Debug(pr.Data[i])
		r, err := sensor.NewReading(pr.Data[i])
		if err != nil {
			bad_payload_handler(w, err)
			return
		}

		err = r.Save(s.db)
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

func init_db() (*sql.DB, error) {
	host, user, password, dbname := "database", "pguser", "pgpassword", "code_challenge"

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	return sql.Open("postgres", psqlInfo)
}
