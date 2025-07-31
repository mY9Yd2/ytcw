package api

import (
	"encoding/json"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog"
	"net/http"
)

type CategoryHandler struct {
	Logger          zerolog.Logger
	CategoryService service.CategoryService
}

func NewCategoryHandler(logger zerolog.Logger, categoryService service.CategoryService) *CategoryHandler {
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
//	@Success		200	{object} model.PaginationResponse[model.CategoryResponse]{data=[]model.CategoryResponse,pagination=model.Pagination}
//	@Router			/categories [get]
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value("pagination").(*model.Pagination)

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
