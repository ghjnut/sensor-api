package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
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

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := &Service{db: db}

	log.Fatal(http.ListenAndServe(":8000", s))

}

type PayloadRaw struct {
	Data []string
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
		s.logHandler(w, r)
	case "/device/":
		s.deviceHandler(w, r)
	}
}

func (s *Service) logHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		badPayloadHandler(w, errors.New("bad request type"))
		return
	}

	dec := json.NewDecoder(req.Body)
	// TODO better way to make this less of a formal struct
	var pr PayloadRaw
	err := dec.Decode(&pr)
	if err != nil {
		badPayloadHandler(w, err)
		return
	}

	// TODO this could fail part-way through. we probably should separate parsing and saving
	// todo for _, data := range pr.Data
	for i := 0; i < len(pr.Data); i++ {
		log.Debug(pr.Data[i])
		r, err := sensor.NewLog(pr.Data[i])
		if err != nil {
			badPayloadHandler(w, err)
			return
		}

		err = r.Save(s.db)
		if err != nil {
			badPayloadHandler(w, err)
			return
		}
	}
}

func (s *Service) deviceHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		badPayloadHandler(w, errors.New("bad request type"))
		return
	}

	device_id := strings.TrimPrefix(req.URL.Path, "/device/")
	_, err := getDeviceLogs(s.db, device_id)
	// TODO correct response for not found, bad payload etc.
	if err != nil {
		badPayloadHandler(w, errors.New("bad request type"))
		return
	}
}

func getDeviceLogs(db *sql.DB, device_id string) ([]sensor.Log, error) {
	rows, err := db.Query("SELECT event_date, temp_farenheit FROM logs WHERE device_id = ?", device_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []sensor.Log

	for rows.Next() {
		var log sensor.Log
		if err := rows.Scan(&log.Date, &log.TempF); err != nil {
			return logs, err
		}
		logs = append(logs, log)
	}
	// redundant, but explicit
	if err = rows.Err(); err != nil {
		return logs, err
	}

	return logs, err
}

func initDB() (*sql.DB, error) {
	host, user, password, dbname := "database", "pguser", "pgpassword", "code_challenge"

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	return sql.Open("postgres", psqlInfo)
}

// what to return if bad decode? should surface, as we would _never_ expect to receive bad payloads
func badPayloadHandler(w http.ResponseWriter, err error) {
	log.Error(err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}
