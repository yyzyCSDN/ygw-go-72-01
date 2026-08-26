package main

import (
	"net/http"
	"path/filepath"

	"powergw/internal/service"
)

func NewServer(gw *service.Gateway, cfg Config) *http.Server {
	handlers := &Handlers{gw: gw, cfg: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("/console", handlers.console)
	mux.HandleFunc("/api/status", handlers.status)
	mux.HandleFunc("/api/points", handlers.points)
	mux.HandleFunc("/api/channels", handlers.channels)
	mux.HandleFunc("/api/versions", handlers.versions)
	mux.HandleFunc("/api/ingest", handlers.ingest)
	mux.HandleFunc("/api/version", handlers.version)
	mux.HandleFunc("/api/timesync", handlers.timesync)
	mux.HandleFunc("/api/rotate", handlers.rotate)
	mux.HandleFunc("/api/flush", handlers.flush)
	mux.HandleFunc("/api/fault", handlers.fault)
	mux.HandleFunc("/api/recover", handlers.recover)
	mux.HandleFunc("/api/demo", handlers.demo)
	return &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
}

func (h *Handlers) console(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(h.cfg.WebDir, "console.html"))
}
