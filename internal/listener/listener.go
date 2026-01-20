package listener

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"tcpeek/internal/executor"
)

type Listener struct {
	IP     string
	Port   int
	Events map[string]string
	conn   net.Conn
}

func New(ip string, port int, events map[string]string) *Listener {
	return &Listener{
		IP:     ip,
		Port:   port,
		Events: events,
	}
}

func (l *Listener) Addr() string {
	return fmt.Sprintf("%s:%d", l.IP, l.Port)
}

func (l *Listener) Start() error {
	conn, err := net.Dial("tcp", l.Addr())
	if err != nil {
		return err
	}
	l.conn = conn

	log.Printf("[INFO] Connected to %s (%d events configured)", l.Addr(), len(l.Events))

	go l.handle()
	return nil
}

func (l *Listener) Stop() error {
	if l.conn != nil {
		return l.conn.Close()
	}
	return nil
}

func (l *Listener) handle() {
	defer l.conn.Close()
	scanner := bufio.NewScanner(l.conn)

	for scanner.Scan() {
		payload := strings.TrimSpace(scanner.Text())
		if payload == "" {
			continue
		}

		log.Printf("[INFO] [%s] Received: %q", l.Addr(), payload)

		cmd, ok := l.match(payload)
		if !ok {
			log.Printf("[WARN] [%s] No match for: %q", l.Addr(), payload)
			continue
		}

		executor.Run(cmd)
	}
}

func (l *Listener) match(payload string) (string, bool) {
	// Try JSON: extract all leaf values and match against event keys
	var data map[string]any
	if err := json.Unmarshal([]byte(payload), &data); err == nil {
		for _, value := range extractLeafValues(data) {
			if cmd, ok := l.Events[value]; ok {
				return strings.ReplaceAll(cmd, "{value}", value), true
			}
		}
	}

	// Fall back to exact string match
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

func extractPath(data map[string]any, path string) (string, bool) {
	parts := strings.Split(path, ".")
	var current any = data

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return "", false
		}
		current, ok = m[part]
		if !ok {
			return "", false
		}
	}

	switch v := current.(type) {
	case string:
		return v, true
	case float64:
		return fmt.Sprintf("%v", v), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}
