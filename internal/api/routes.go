package api

import (
	"github.com/go-chi/chi/v5"
	api "github.com/mY9Yd2/ytcw/internal/api/handlers"
	"github.com/mY9Yd2/ytcw/internal/service"
	"github.com/rs/zerolog"
)

func Routes(logger zerolog.Logger,
	channelService service.ChannelService,
	videoService service.VideoService,
	categoryService service.CategoryService) chi.Router {
	r := chi.NewRouter()

	channelHandler := api.NewChannelHandler(logger, channelService)
	videoHandler := api.NewVideoHandler(logger, videoService)
	categoryHandler := api.NewCategoryHandler(logger, categoryService)

	r.Route("/channels", func(r chi.Router) {
		r.With(PaginationMiddleware).Get("/", channelHandler.ListChannels)
	})

	r.Route("/videos", func(r chi.Router) {
		r.With(PaginationMiddleware).Get("/", videoHandler.ListVideos)
	})

	r.Route("/categories", func(r chi.Router) {
		r.With(PaginationMiddleware).Get("/", categoryHandler.ListCategories)
	})

	return r
}
