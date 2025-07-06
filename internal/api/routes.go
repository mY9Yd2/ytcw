package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"ytcw/internal/service"
)

func Routes(logger zerolog.Logger, channelService service.ChannelService) chi.Router {
	r := chi.NewRouter()

	channelHandler := NewChannelHandler(logger, channelService)

	r.Get("/channels", channelHandler.ListChannels)

	return r
}
