package content

import (
	"encoding/json"
	"net/http"

	"github.com/mY9Yd2/ytcw/internal/common"
	"github.com/rs/zerolog"
)

type CategoryHandler struct {
	Logger          zerolog.Logger
	CategoryService CategoryService
}

func NewCategoryHandler(logger zerolog.Logger, categoryService CategoryService) *CategoryHandler {
	return &CategoryHandler{
		Logger:          logger,
		CategoryService: categoryService,
	}
}

// ListCategories godoc
//
//	@Summary		List categories
//	@Description	Get a paginated list of categories ordered by name in ascending order
//	@Tags			Categories
//	@Produce		json
//	@Param			page		query int false "page"
//	@Param			page_size	query int false "page size"
//	@Success		200	{object} common.PaginationResponse[CategoryResponse]{data=[]CategoryResponse,pagination=common.Pagination}
//	@Router			/categories [get]
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value("pagination").(*common.Pagination)

	categories, pageMeta, err := h.CategoryService.GetCategories(p)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to get categories")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := common.PaginationResponse[CategoryResponse]{
		Data:       categories,
		Pagination: pageMeta,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
