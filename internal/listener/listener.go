package listener

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"tcpeek/internal/executor"
)

const maxBackoff = 30 * time.Second

type Listener struct {
	IP            string
	Port          int
	Events        map[string]string
	debug         bool
	autoReconnect bool
	mu            sync.Mutex
	conn          net.Conn
	done          chan struct{}
	reconnect     chan struct{}
}

func New(ip string, port int, events map[string]string, debug, autoReconnect bool) *Listener {
	return &Listener{
		IP:            ip,
		Port:          port,
		Events:        events,
		debug:         debug,
		autoReconnect: autoReconnect,
		done:          make(chan struct{}),
		reconnect:     make(chan struct{}, 1),
	}
}

func (l *Listener) Reconnect() {
	select {
	case l.reconnect <- struct{}{}:
	default:
	}
}

func (l *Listener) Addr() string {
	return fmt.Sprintf("%s:%d", l.IP, l.Port)
}

func (l *Listener) Start() error {
	go l.run()
	return nil
}

func (l *Listener) Stop() error {
	close(l.done)
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.conn != nil {
		return l.conn.Close()
	}
	return nil
}

func (l *Listener) run() {
	backoff := time.Second

	for {
		select {
		case <-l.done:
			return
		default:
		}

		conn, err := net.Dial("tcp", l.Addr())
		if err != nil {
			if !l.autoReconnect {
				if l.debug {
					log.Printf("[DEBUG] [%s] Connect failed, waiting for reconnect signal", l.Addr())
				}
				select {
				case <-l.done:
					return
				case <-l.reconnect:
				}
				continue
			}
			if l.debug {
				log.Printf("[DEBUG] [%s] Connect failed, retrying in %s", l.Addr(), backoff)
			}
			select {
			case <-l.done:
				return
			case <-l.reconnect:
				backoff = time.Second
			case <-time.After(backoff):
				backoff = min(backoff*2, maxBackoff)
			}
			continue
		}

		backoff = time.Second
		log.Printf("[INFO] Connected to %s (%d events configured)", l.Addr(), len(l.Events))

		l.mu.Lock()
		l.conn = conn
		l.mu.Unlock()

		l.handle(conn)

		log.Printf("[INFO] [%s] Connection lost, reconnecting", l.Addr())
	}
}

func (l *Listener) handle(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		payload := strings.TrimSpace(scanner.Text())
		if payload == "" {
			continue
		}

		if l.debug {
			log.Printf("[DEBUG] [%s] Received: %q", l.Addr(), payload)
		}

		cmd, ok := l.match(payload)
		if !ok {
			log.Printf("[WARN] [%s] No match for: %q", l.Addr(), payload)
			continue
		}

		if l.debug {
			log.Printf("[DEBUG] [%s] Event %q → %s", l.Addr(), payload, cmd)
		}
		executor.Run(cmd)
	}
}

func (l *Listener) match(payload string) (string, bool) {
	var data map[string]any
	if err := json.Unmarshal([]byte(payload), &data); err == nil {
		for _, value := range extractLeafValues(data) {
			if cmd, ok := l.Events[value]; ok {
				return strings.ReplaceAll(cmd, "{value}", value), true
			}
		}
	}

	cmd, ok := l.Events[payload]
	return cmd, ok
}

func extractLeafValues(data map[string]any) []string {
	var values []string
	for _, v := range data {
		switch val := v.(type) {
		case string:
			values = append(values, val)
		case map[string]any:
			values = append(values, extractLeafValues(val)...)
		}
	}
	return values
}
