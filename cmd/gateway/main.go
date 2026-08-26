package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"powergw/internal/service"
)

func main() {
	addr := flag.String("addr", "", "listen address")
	configPath := flag.String("config", "", "config file path")
	dir := flag.String("dir", "", "data directory")
	webDir := flag.String("web", "", "web directory")
	flag.Parse()

	cfg := DefaultConfig()
	if *configPath != "" {
		loaded, err := LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("load config: %v", err)
		}
		cfg = cfg.Merge(loaded)
	}
	cfg = cfg.Merge(Config{
		Addr:    *addr,
		DataDir: *dir,
		WebDir:  *webDir,
	})
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	gw := BuildGateway()
	runner := service.NewRunner(gw, time.Duration(cfg.CycleSeconds)*time.Second)
	runner.Start()

	server := NewServer(gw, cfg)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
		runner.Stop()
		gw.CloseSessions()
	}()

	log.Printf("power protocol gateway listening on %s", cfg.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
