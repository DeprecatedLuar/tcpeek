package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tcpeek/internal/config"
	"tcpeek/internal/listener"
)

func start() {
	log.Println("[INFO] tcpeek starting")

	if err := writePID(); err != nil {
		log.Fatalf("[ERROR] Failed to write PID file: %v", err)
	}
	defer removePID()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[ERROR] Failed to load config: %v", err)
	}

	if len(cfg.Listeners) == 0 {
		log.Fatalf("[ERROR] No listeners configured in %s", config.ConfigDir)
	}

	var listeners []*listener.Listener
	for _, lc := range cfg.Listeners {
		l := listener.New(lc.IP, lc.Port, lc.Events)
		if err := l.Start(); err != nil {
			log.Printf("[ERROR] Failed to start listener %s:%d: %v", lc.IP, lc.Port, err)
			continue
		}
		listeners = append(listeners, l)
	}

	if len(listeners) == 0 {
		log.Fatal("[ERROR] No listeners started")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("[INFO] Shutting down")
	for _, l := range listeners {
		l.Stop()
	}
}
