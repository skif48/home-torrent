package config

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/samber/do"
)

var (
	k      = koanf.New(".")
	parser = yaml.Parser()
)

type LoggingConfig struct {
	LogLevel string `koanf:"log_level"`
}

type TorrentConfig struct {
	PeerId   [20]byte `koanf:"peer_id"`
	PeerPort uint16   `koanf:"peer_port"`
}

type Config struct {
	HttpPort int            `koanf:"http_port"`
	Logging  *LoggingConfig `koanf:"logging"`
	Torrent  *TorrentConfig `koanf:"torrent"`
}

// for some reason doesn't take into account config file content at all
func NewConfig(i *do.Injector) (*Config, error) {
	path := os.Getenv(`CONFIG_PATH`)
	if path == "" {
		path = "./home_torrent_config.yaml"
	}
	if err := k.Load(file.Provider(path), parser); err != nil {
		return nil, err
	}

	config := Config{}

	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
