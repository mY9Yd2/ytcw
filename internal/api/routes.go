package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog"
)

func Routes(logger zerolog.Logger,
	channelService service.ChannelService,
	videoService service.VideoService) chi.Router {
	r := chi.NewRouter()

	channelHandler := NewChannelHandler(logger, channelService, videoService)

	r.Get("/channels", channelHandler.ListChannels)
	r.Get("/videos", channelHandler.ListVideos)

	return r
}
