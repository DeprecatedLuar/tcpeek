package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DeprecatedLuar/tcpeek/internal/config"
	"github.com/DeprecatedLuar/tcpeek/internal/listener"
	"github.com/DeprecatedLuar/tcpeek/internal/logfile"
)

func start(debug bool) {
	if w, err := logfile.Open(); err != nil {
		log.Printf("[WARN] file logging disabled: %v", err)
	} else {
		defer w.Close()
		log.SetOutput(io.MultiWriter(os.Stderr, w))
	}

	if err := daemon.Run(func(ctx context.Context) error { return serve(ctx, debug) }); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}

func serve(ctx context.Context, debug bool) error {
	log.Println("[INFO] tcpeek starting")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Listeners) == 0 {
		return fmt.Errorf("no listeners configured in %s", config.ConfigDir)
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

	<-ctx.Done()

	log.Println("[INFO] Shutting down")
	for _, l := range listeners {
		l.Stop()
	}
	return nil
}
