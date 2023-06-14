package httptransport

import (
	"encoding/json"
	"net/http"
	"strings"

	"ghjnut/sensor/internal"
)

type handler struct {
	internal.Service
}

func Handler(s internal.Service) http.Handler {
	h := &handler{s}
	r := http.NewServeMux()

	//case r.URL.Path == "/ingest":
	//	s.logHandler(w, r)
	//case strings.HasPrefix(r.URL.Path, "/device/"):
	//	s.deviceHandler(w, r)

	r.HandleFunc("/ingest", h.createLogs)
	r.HandleFunc("/device/", h.device)
	return r
}

func (h *handler) createLogs(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	dec := json.NewDecoder(req.Body)
	// TODO better way to make this less of a formal struct
	var in internal.CreateLogsIn
	defer req.Body.Close()
	err := dec.Decode(&in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	out, err := h.CreateLogs(ctx, in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (h *handler) device(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	device_id := strings.TrimPrefix(req.URL.Path, "/device/")
	var in internal.DeviceIn
	in.ID = device_id

	ctx := req.Context()

	out, err := h.Device(ctx, in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
