package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog"
)

func Routes(logger zerolog.Logger, channelService service.ChannelService) chi.Router {
	r := chi.NewRouter()

	channelHandler := NewChannelHandler(logger, channelService)

	r.Get("/channels", channelHandler.ListChannels)

	return r
}
