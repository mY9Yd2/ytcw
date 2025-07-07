package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var Pretty zerolog.Logger
var JSON zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	Pretty = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Timestamp().Logger()

	JSON = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Set the default logger
	log.Logger = JSON
}
