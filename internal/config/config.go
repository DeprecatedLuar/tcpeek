package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	ConfigDirName = "tcpeek"
	TomlExt       = ".toml"
	EnvConfigDir  = "TCPEEK_CONFIG_DIR"
)

var (
	ConfigDir string
)

func init() {
	if dir := os.Getenv(EnvConfigDir); dir != "" {
		ConfigDir = dir
		return
	}

	configBase := os.Getenv("XDG_CONFIG_HOME")
	if configBase == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("[ERROR] resolving home directory: %v", err)
		}
		configBase = filepath.Join(home, ".config")
	}
	ConfigDir = filepath.Join(configBase, ConfigDirName)
}

type Config struct {
	Listeners []Listener
}

type Listener struct {
	IP        string
	Port      int
	Events    map[string]string
	Reconnect bool
}

type tomlConfig struct {
	Reconnect *bool             `toml:"reconnect"`
	Events    map[string]string `toml:"events"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	ipDirs, err := os.ReadDir(ConfigDir)
	if err != nil {
		return nil, fmt.Errorf("reading config dir: %w", err)
	}

	for _, ipDir := range ipDirs {
		if !ipDir.IsDir() {
			continue
		}
		ip := ipDir.Name()

		ipPath := filepath.Join(ConfigDir, ip)
		portFiles, err := os.ReadDir(ipPath)
		if err != nil {
			log.Printf("[WARN] skipping %s: %v", ipPath, err)
			continue
		}

		for _, portFile := range portFiles {
			if portFile.IsDir() || !strings.HasSuffix(portFile.Name(), TomlExt) {
				continue
			}

			filePath := filepath.Join(ipPath, portFile.Name())
			portStr := strings.TrimSuffix(portFile.Name(), TomlExt)
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Printf("[WARN] skipping %s: invalid port", filePath)
				continue
			}

			listener, err := parseFile(filePath, ip, port)
			if err != nil {
				log.Printf("[WARN] skipping %s: %v", filePath, err)
				continue
			}

			cfg.Listeners = append(cfg.Listeners, *listener)
		}
	}

	return cfg, nil
}

func parseFile(path, ip string, port int) (*Listener, error) {
	var tc tomlConfig
	if _, err := toml.DecodeFile(path, &tc); err != nil {
		return nil, err
	}

	reconnect := tc.Reconnect == nil || *tc.Reconnect

	return &Listener{
		IP:        ip,
		Port:      port,
		Events:    tc.Events,
		Reconnect: reconnect,
	}, nil
}

