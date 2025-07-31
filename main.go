package main

import (
	"github.com/mY9Yd2/ytcw/cmd"
	_ "github.com/mY9Yd2/ytcw/docs"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command")
	}
}
