package api

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
	model "ytcw/internal/model/api"
	"ytcw/internal/service"
)

type ChannelHandler struct {
	Logger         zerolog.Logger
	ChannelService service.ChannelService
}

func NewChannelHandler(logger zerolog.Logger, channelService service.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		Logger:         logger,
		ChannelService: channelService,
	}
}

func (h *ChannelHandler) ListChannels(w http.ResponseWriter, r *http.Request) {
	p := model.NewPaginationFromRequest(r)

	channels, pageMeta, err := h.ChannelService.GetChannels(p)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get channels")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"data":       channels,
		"pagination": pageMeta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error().Err(err).Msg("failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
