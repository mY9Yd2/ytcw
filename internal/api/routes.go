package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mY9Yd2/ytcw/internal/content"
	"github.com/rs/zerolog"
)

func Routes(logger zerolog.Logger,
	channelService content.ChannelService,
	videoService content.VideoService,
	categoryService content.CategoryService) chi.Router {
	r := chi.NewRouter()

	r.Mount("/channels", content.ChannelRoutes(logger, channelService))
	r.Mount("/videos", content.VideoRoutes(logger, videoService))
	r.Mount("/categories", content.CategoryRoutes(logger, categoryService))

	return r
}
