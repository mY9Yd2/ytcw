package cmd

import (
	"github.com/go-chi/chi/v5"
	"github.com/mY9Yd2/ytcw/internal/api"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/logger"
	"github.com/mY9Yd2/ytcw/internal/repository"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/spf13/cobra"
	"net/http"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Run:     serve,
	GroupID: "main",
}

func serve(cmd *cobra.Command, args []string) {
	log := logger.JSON

	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	channelRepo := repository.NewChannelRepository(dbCon)
	channelService := service.NewChannelService(channelRepo)

	r := chi.NewRouter()
	r.Use(api.ZerologMiddleware(log))

	r.Mount("/api/v1", api.Routes(log, channelService))

	log.Info().Msg("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal().Err(err).Msg("server exited with error")
	}
}
