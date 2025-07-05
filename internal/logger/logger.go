package logger

import (
	"github.com/rs/zerolog"
	"os"
)

var Pretty zerolog.Logger
var JSON zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	Pretty = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Timestamp().Logger()

	JSON = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
