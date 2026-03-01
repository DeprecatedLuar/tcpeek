package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	ConfigDirName = "tcpeek"
	TomlExt       = ".toml"
)

var (
	ConfigDir string
)

func init() {
	configBase := os.Getenv("XDG_CONFIG_HOME")
	if configBase == "" {
		home, _ := os.UserHomeDir()
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

		portFiles, err := os.ReadDir(filepath.Join(ConfigDir, ip))
		if err != nil {
			continue
		}

		for _, portFile := range portFiles {
			if portFile.IsDir() || !strings.HasSuffix(portFile.Name(), TomlExt) {
				continue
			}

			portStr := strings.TrimSuffix(portFile.Name(), TomlExt)
			port, err := strconv.Atoi(portStr)
			if err != nil {
				continue
			}

			filePath := filepath.Join(ConfigDir, ip, portFile.Name())
			listener, err := parseFile(filePath, ip, port)
			if err != nil {
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

