package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ServiceLogger wraps zerolog with service-specific context
type ServiceLogger struct {
	logger zerolog.Logger
}

// Config represents logger configuration
type Config struct {
	Level   string `json:"level"`
	Pretty  bool   `json:"pretty"`
	Service string `json:"service"`
}

// New creates a new ServiceLogger instance
func New(config Config) *ServiceLogger {
	// Set global log level
	level, err := zerolog.ParseLevel(strings.ToLower(config.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var logger zerolog.Logger

	if config.Pretty {
		// Use console writer for pretty output with custom formatting
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				level := strings.ToUpper(i.(string))
				switch level {
				case "INFO":
					return fmt.Sprintf("\x1b[32m[%s]\x1b[0m", level) // Green
				case "WARN":
					return fmt.Sprintf("\x1b[33m[%s]\x1b[0m", level) // Yellow
				case "ERROR":
					return fmt.Sprintf("\x1b[31m[%s]\x1b[0m", level) // Red
				case "DEBUG":
					return fmt.Sprintf("\x1b[35m[%s]\x1b[0m", level) // Magenta.. not sure why i chose this lol
				case "FATAL":
					return fmt.Sprintf("\x1b[41m[%s]\x1b[0m", level) // Red Background
				default:
					return fmt.Sprintf("[%s]", level)
				}
			},
			FormatMessage: func(i interface{}) string { return fmt.Sprintf("%s", i) },
		}
		logger = zerolog.New(consoleWriter).With().
			Timestamp().
			Str("service", config.Service).
			Logger()
	} else {
		// Use default JSON output to stdout
		logger = zerolog.New(os.Stdout).With().
			Timestamp().
			Str("service", config.Service).
			Logger()
	}

	return &ServiceLogger{
		logger: logger,
	}
}

// WithField adds a field to the logger context
func (sl *ServiceLogger) WithField(key string, value interface{}) *ServiceLogger {
	return &ServiceLogger{
		logger: sl.logger.With().Interface(key, value).Logger(),
	}
}

// WithFields adds multiple fields to the logger context
func (sl *ServiceLogger) WithFields(fields map[string]interface{}) *ServiceLogger {
	ctx := sl.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &ServiceLogger{
		logger: ctx.Logger(),
	}
}

// WithError adds an error to the logger context
func (sl *ServiceLogger) WithError(err error) *ServiceLogger {
	return &ServiceLogger{
		logger: sl.logger.With().Err(err).Logger(),
	}
}

// Debug logs a debug message
func (sl *ServiceLogger) Debug(msg string) {
	sl.logger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func (sl *ServiceLogger) Debugf(format string, args ...interface{}) {
	sl.logger.Debug().Msgf(format, args...)
}

// Info logs an info message
func (sl *ServiceLogger) Info(msg string) {
	sl.logger.Info().Msg(msg)
}

// Infof logs a formatted info message
func (sl *ServiceLogger) Infof(format string, args ...interface{}) {
	sl.logger.Info().Msgf(format, args...)
}

// Warn logs a warning message
func (sl *ServiceLogger) Warn(msg string) {
	sl.logger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func (sl *ServiceLogger) Warnf(format string, args ...interface{}) {
	sl.logger.Warn().Msgf(format, args...)
}

// Error logs an error message
func (sl *ServiceLogger) Error(msg string) {
	sl.logger.Error().Msg(msg)
}

// Errorf logs a formatted error message
func (sl *ServiceLogger) Errorf(format string, args ...interface{}) {
	sl.logger.Error().Msgf(format, args...)
}

// Fatal logs a fatal message and exits
func (sl *ServiceLogger) Fatal(msg string) {
	sl.logger.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal message and exits
func (sl *ServiceLogger) Fatalf(format string, args ...interface{}) {
	sl.logger.Fatal().Msgf(format, args...)
}

// GetZerologLogger returns the underlying zerolog.Logger for advanced usage
func (sl *ServiceLogger) GetZerologLogger() zerolog.Logger {
	return sl.logger
}

// SetGlobalLogger sets the global zerolog logger
func SetGlobalLogger(config Config) {
	logger := New(config)
	log.Logger = logger.logger
}
