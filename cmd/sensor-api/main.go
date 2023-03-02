package main

import (
	"encoding/json"
	"errors"
	"io"
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

func ingestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		bad_payload_handler(w, errors.New("bad request type 'POST'"))
		return
	}

	dec := json.NewDecoder(req.Body)
	var sp sensor.SensorPayload
	err := dec.Decode(&sp)
	// seems odd to put EOF in the class of "errors"
	if err != io.EOF && err != nil {
		bad_payload_handler(w, err)
		return
	}

	log.Debug(sp.Data)

	err = sp.Validate()
	if err != nil {
		bad_payload_handler(w, err)
		return
	}

	err = sp.Save()
	if err != nil {
		bad_payload_handler(w, err)
		return
	}
}

// what to return if bad decode? should surface, as we would _never_ expect to receive bad payloads
func bad_payload_handler(w http.ResponseWriter, err error) {
	log.Error(err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}
