package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tcpeek/internal/config"
	"tcpeek/internal/listener"
)

func start(debug bool) {
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
		l := listener.New(lc.IP, lc.Port, lc.Events, debug, lc.Reconnect)
		l.Start()
		listeners = append(listeners, l)
	}

	if debug {
		log.Printf("[DEBUG] %d listener(s) configured:", len(listeners))
		for _, l := range listeners {
			log.Printf("[DEBUG]   %s (%d events)", l.Addr(), len(l.Events))
		}
	}

	usr1 := make(chan os.Signal, 1)
	signal.Notify(usr1, syscall.SIGUSR1)
	go func() {
		for range usr1 {
			log.Println("[INFO] Reconnecting all listeners")
			for _, l := range listeners {
				l.Reconnect()
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("[INFO] Shutting down")
	for _, l := range listeners {
		l.Stop()
	}
}
