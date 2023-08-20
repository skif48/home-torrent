package config

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type LoggingConfig struct {
	LogLevel string `koanf:"log_level"`
}

type TorrentConfig struct {
	PeerId   [20]byte
	PeerPort uint16
}

type Config struct {
	HttpPort int            `koanf:"http_port"`
	Logging  *LoggingConfig `koanf:"logging"`
	Torrent  *TorrentConfig `koanf:"torrent"`
}

func InitConfig() (error, *Config) {
	path := os.Getenv(`CONFIG_PATH`)
	if path == "" {
		path = "./torrent_config.yaml"
	}
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return err, nil
	}

	config := &Config{}

	if err := k.UnmarshalWithConf(".", config, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return err, nil
	}

	return nil, config
}
