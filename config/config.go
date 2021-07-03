package config

import (
	_ "embed"
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"

	"github.com/xeipuuv/gojsonschema"
	"vladusenko.io/home-torrent/defaults"
)

type LoggingConfig struct {
	LogLevel   string `json:"log_level"`
	Console    bool   `json:"console"`
	Directory  string `json:"directory"`
	Filename   string `json:"file_name"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
}
type Config struct {
	HttpPort int            `json:"http_port"`
	Logging  *LoggingConfig `json:"logging"`
}

//go:embed schema.json
var configSchema string

var config *Config = nil
var schema *gojsonschema.Schema = nil
var once *sync.Once = new(sync.Once)

// NOTE for unit tests only
func Reset() {
	once = new(sync.Once)
}

func GetConfig(path string) (*Config, error) {
	var err error
	once.Do(func() {
		schemaLoader := gojsonschema.NewStringLoader(configSchema)
		if schema, err = gojsonschema.NewSchema(schemaLoader); err != nil {
			panic(err)
		}
		config, err = parseAndValidateConfig(path)
	})

	return config, err
}

func parseAndValidateConfig(path string) (*Config, error) {
	var raw []byte
	var validationResult *gojsonschema.Result
	var err error

	if raw, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	loader := gojsonschema.NewBytesLoader(raw)

	if validationResult, err = schema.Validate(loader); err != nil {
		return nil, err
	}

	if !validationResult.Valid() {
		return nil, errors.New("config file is invalid")
	}

	config := &Config{
		HttpPort: defaults.DEFAULT_HTTP_PORT,
		Logging: &LoggingConfig{
			LogLevel:   "info",
			Console:    true,
			Directory:  "./logs",
			Filename:   "ht.log",
			MaxSize:    10,
			MaxBackups: 25,
			MaxAge:     30,
		},
	}

	if err = json.Unmarshal(raw, config); err != nil {
		return nil, err
	}

	return config, nil
}
