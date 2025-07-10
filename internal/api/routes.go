package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog"
)

func Routes(logger zerolog.Logger,
	channelService service.ChannelService,
	videoService service.VideoService,
	categoryService service.CategoryService) chi.Router {
	r := chi.NewRouter()

	channelHandler := NewChannelHandler(logger, channelService, videoService, categoryService)

	r.Get("/channels", channelHandler.ListChannels)
	r.Get("/videos", channelHandler.ListVideos)
	r.Get("/categories", channelHandler.ListCategories)

	return r
}
