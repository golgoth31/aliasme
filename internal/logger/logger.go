package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Configure sets up the logger based on configuration
func Configure() {
	// Set log level
	level := viper.GetString("logging.level")
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Configure time format
	timeFormat := viper.GetString("logging.time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	zerolog.TimeFieldFormat = timeFormat

	// Configure output format
	format := viper.GetString("logging.format")
	if format == "json" {
		// Configure ECS format
		zerolog.TimestampFieldName = "@timestamp"
		zerolog.LevelFieldName = "log.level"
		zerolog.MessageFieldName = "message"
		zerolog.ErrorFieldName = "error.message"
		zerolog.ErrorStackFieldName = "error.stack_trace"

		// Create logger with ECS fields
		log.Logger = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Str("ecs.version", "1.6.0").
			Str("service.name", "aliasme").
			Str("service.type", "grpc").
			Logger()
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: timeFormat,
		}
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}
}
