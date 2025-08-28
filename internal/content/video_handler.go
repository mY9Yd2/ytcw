package content

import (
	"encoding/json"
	"net/http"

	"github.com/mY9Yd2/ytcw/internal/common"
	"github.com/rs/zerolog"
)

type VideoHandler struct {
	Logger       zerolog.Logger
	VideoService VideoService
}

func NewVideoHandler(logger zerolog.Logger, videoService VideoService) *VideoHandler {
	return &VideoHandler{
		Logger:       logger,
		VideoService: videoService,
	}
}

// ListVideos godoc
//
//	@Summary		List videos
//	@Description	Get a paginated list of videos ordered by created_at in descending order
//	@Tags			Videos
//	@Produce		json
//	@Param			page		query int false "page"
//	@Param			page_size	query int false "page size"
//	@Success		200	{object} common.PaginationResponse[VideoResponse]{data=[]VideoResponse,pagination=common.Pagination}
//	@Router			/videos [get]
func (h *VideoHandler) ListVideos(w http.ResponseWriter, r *http.Request) {
	p := common.NewPaginationFromRequest(r)

	videos, pageMeta, err := h.VideoService.GetVideos(p)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to get videos")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := common.PaginationResponse[VideoResponse]{
		Data:       videos,
		Pagination: pageMeta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
