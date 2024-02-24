package routers

import (
	"github.com/GabDewraj/library-api/pkgs/api/middleware"
	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

type LibraryRouterParams struct {
	fx.In
	Mux        *chi.Mux
	Middleware middleware.Service
	Handler    books.Handler
}

func NewBooksRouter(params LibraryRouterParams) {
	// Create a fresh middleware stack so that other domains

	params.Mux.Route("/books", func(r chi.Router) {
		// Logging
		r.Use(params.Middleware.CustomLogger)
		// Add CORS for browsers
		r.Use(params.Middleware.CORS)
		// Add rate limiting
		r.Use(params.Middleware.RateLimiter)
		// Routes
		r.Post("/", params.Handler.CreateBook)
		r.Get("/", params.Handler.GetBooks)
		r.Get("/{book_id}", params.Handler.GetBookByID)
		r.Put("/{book_id}", params.Handler.UpdateBook)
		r.Delete("/{book_id}", params.Handler.DeleteBook)
	})

}
