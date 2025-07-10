package api

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	model "github.com/mY9Yd2/ytcw/internal/model/api"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

func ZerologMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now().UTC()

			next.ServeHTTP(ww, r)

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("duration", time.Since(t1)).
				Msg("http request")
		})
	}
}

func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := model.NewPaginationFromRequest(r)
		ctx := context.WithValue(r.Context(), "pagination", p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
