package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"powergw/internal/model"
	"powergw/internal/service"
)

type Handlers struct {
	gw  *service.Gateway
	cfg Config
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func (h *Handlers) status(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.gw.Status())
}

func (h *Handlers) points(w http.ResponseWriter, r *http.Request) {
	snap := h.gw.Mapper.Snapshot()
	writeJSON(w, http.StatusOK, map[string]any{
		"table":  snap.TableID,
		"points": snap.Points,
	})
}

func (h *Handlers) channels(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.gw.Channels.Snapshot())
}

func (h *Handlers) versions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.gw.Controller.Snapshot())
}

func (h *Handlers) ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	var body struct {
		Channel string `json:"channel"`
		Hex     string `json:"hex"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	raw, err := hex.DecodeString(body.Hex)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.gw.Ingest(body.Channel, raw); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	var body struct {
		ID       string `json:"id"`
		Proto    string `json:"proto"`
		Table    string `json:"table"`
		Checksum uint64 `json:"checksum"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	proto, ok := model.ProtocolFromString(body.Proto)
	if !ok {
		writeError(w, http.StatusBadRequest, fmt.Errorf("unknown protocol"))
		return
	}
	if err := h.gw.ApplyVersion(model.NewProtocolVersion(body.ID, proto, body.Table, body.Checksum)); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) timesync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	count, err := h.gw.SyncAll()
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"synced": count})
}

func (h *Handlers) rotate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	count, err := h.gw.RotateAll()
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"rotated": count})
}

func (h *Handlers) flush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	count, err := h.gw.FlushAll()
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"flushed": count})
}

func (h *Handlers) fault(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	var body struct {
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.gw.FaultChannel(body.Channel); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) recover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	if err := h.gw.Recover(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) demo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	var body struct {
		Rounds int `json:"rounds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if body.Rounds <= 0 {
		body.Rounds = 1
	}
	count, err := service.FeedDemo(h.gw, body.Rounds)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"ingested": count})
}
