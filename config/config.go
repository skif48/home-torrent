package config

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/samber/do"
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

func NewConfig(i *do.Injector) (*Config, error) {
	path := os.Getenv(`CONFIG_PATH`)
	if path == "" {
		path = "./torrent_config.yaml"
	}
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, err
	}

	config := &Config{}

	if err := k.UnmarshalWithConf(".", config, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, err
	}

	return config, nil
}
