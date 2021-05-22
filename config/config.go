package config

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpPort int `yaml:"http_port"`
}

var config *Config = nil
var once sync.Once

func GetConfig(path string) (*Config, error) {
	var err error
	once.Do(func() {
		config, err = parseConfig(path)
	})

	return config, err
}

func parseConfig(path string) (*Config, error) {
	var raw []byte
	var err error

	if raw, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	config := &Config{}

	if err = yaml.Unmarshal(raw, config); err != nil {
		return nil, err
	}

	return config, nil
}
