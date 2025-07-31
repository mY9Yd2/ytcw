package api

import (
	"encoding/json"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog"
	"net/http"
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

// ListChannels godoc
//
//	@Summary		List channels
//	@Description	Get a paginated list of channels ordered by channel in ascending order
//	@Tags			Channels
//	@Produce		json
//	@Param			page		query int false "page"
//	@Param			page_size	query int false "page size"
//	@Param			category	query string false "category name"
//	@Success		200	{object} model.PaginationResponse[model.ChannelResponse]{data=[]model.ChannelResponse,pagination=model.Pagination}
//	@Router			/channels [get]
func (h *ChannelHandler) ListChannels(w http.ResponseWriter, r *http.Request) {
	p := model.NewPaginationFromRequest(r)
	category := r.URL.Query().Get("category")

	channels, pageMeta, err := h.ChannelService.GetChannels(p, category)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to get channels")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := model.PaginationResponse[model.ChannelResponse]{
		Data:       channels,
		Pagination: pageMeta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
