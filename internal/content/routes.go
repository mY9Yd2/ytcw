package content

import (
	"github.com/go-chi/chi/v5"
	"github.com/mY9Yd2/ytcw/internal/common"
	"github.com/rs/zerolog"
)

func ChannelRoutes(logger zerolog.Logger, service ChannelService) chi.Router {
	r := chi.NewRouter()
	handler := NewChannelHandler(logger, service)

	r.With(common.PaginationMiddleware).Get("/", handler.ListChannels)

	return r
}

func VideoRoutes(logger zerolog.Logger, service VideoService) chi.Router {
	r := chi.NewRouter()
	handler := NewVideoHandler(logger, service)

	r.With(common.PaginationMiddleware).Get("/", handler.ListVideos)

	return r
}

func CategoryRoutes(logger zerolog.Logger, service CategoryService) chi.Router {
	r := chi.NewRouter()
	handler := NewCategoryHandler(logger, service)

	r.With(common.PaginationMiddleware).Get("/", handler.ListCategories)

	return r
}
