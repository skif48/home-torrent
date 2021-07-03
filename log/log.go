package log

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"vladusenko.io/home-torrent/config"
)

type Logger struct {
	logger zerolog.Logger
}

func (logger *Logger) Info() *zerolog.Event {
	return logger.logger.Info()
}

var logger *Logger

func newRollingFile(config *config.LoggingConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}

func Configure(config *config.LoggingConfig) {
	var writers []io.Writer

	if config.Console {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	writers = append(writers, newRollingFile(config))

	mw := io.MultiWriter(writers...)

	var logLevel zerolog.Level

	switch config.LogLevel {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		panic(errors.New("Unknown log level: " + config.LogLevel))
	}

	zerolog.SetGlobalLevel(logLevel)
	tempLogger := zerolog.New(mw).With().Timestamp().Logger()
	logger = &Logger{
		logger: tempLogger,
	}

	tempLogger.Info().Msg("Welcome to home-torrent")

	tempLogger.Info().
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("Logging configured")
}

func GetLogger() *Logger {
	if logger != nil {
		return logger
	} else {
		panic(errors.New("Logger is not configured"))
	}
}
