package api

import (
	"encoding/json"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog"
	"net/http"
)

type ChannelHandler struct {
	Logger          zerolog.Logger
	ChannelService  service.ChannelService
	VideoService    service.VideoService
	CategoryService service.CategoryService
}

func NewChannelHandler(logger zerolog.Logger,
	channelService service.ChannelService,
	videoService service.VideoService,
	categoryService service.CategoryService) *ChannelHandler {
	return &ChannelHandler{
		Logger:          logger,
		ChannelService:  channelService,
		VideoService:    videoService,
		CategoryService: categoryService,
	}
}

func (h *ChannelHandler) ListChannels(w http.ResponseWriter, r *http.Request) {
	p := model.NewPaginationFromRequest(r)

	channels, pageMeta, err := h.ChannelService.GetChannels(p)
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

func (h *ChannelHandler) ListVideos(w http.ResponseWriter, r *http.Request) {
	p := model.NewPaginationFromRequest(r)

	videos, pageMeta, err := h.VideoService.GetVideos(p)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to get videos")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := model.PaginationResponse[model.VideoResponse]{
		Data:       videos,
		Pagination: pageMeta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *ChannelHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	p := model.NewPaginationFromRequest(r)

	categories, pageMeta, err := h.CategoryService.GetCategories(p)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to get categories")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := model.PaginationResponse[model.CategoryResponse]{
		Data:       categories,
		Pagination: pageMeta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
