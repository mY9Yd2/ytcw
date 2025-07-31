package cmd

import (
	"github.com/go-chi/chi/v5"
	"github.com/mY9Yd2/ytcw/internal/api"
	"github.com/mY9Yd2/ytcw/internal/config"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/repository"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Start the REST API server",
	Run:     serve,
	GroupID: "main",
}

// @title			ytcw API
// @version			1.0
// @license.name	MIT License
// @license.url		https://github.com/mY9Yd2/ytcw/blob/main/LICENSE.md
// @host			localhost:8080
// @basePath		/api/v1
func serve(cmd *cobra.Command, args []string) {
	cfg := config.GetConfig()
	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	channelRepo := repository.NewChannelRepository(dbCon)
	channelService := service.NewChannelService(channelRepo)
	videoRepo := repository.NewVideoRepository(dbCon)
	videoService := service.NewVideoService(videoRepo)
	categoryRepo := repository.NewCategoryRepository(dbCon)
	categoryService := service.NewCategoryService(categoryRepo)

	r := chi.NewRouter()
	r.Use(api.ZerologMiddleware(log.Logger))

	if cfg.IsDevelopment() {
		r.Mount("/swagger", httpSwagger.WrapHandler)
	}

	r.Mount("/api/v1", api.Routes(log.Logger, channelService, videoService, categoryService))

	log.Info().Str("address", cfg.Api.Address).Msg("Starting server")
	if err := http.ListenAndServe(cfg.Api.Address, r); err != nil {
		log.Fatal().Err(err).Msg("Server exited with error")
	}
}
