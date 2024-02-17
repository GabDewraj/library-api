package routers

import (
	"github.com/GabDewraj/library-api/pkgs/api/handlers"
	"github.com/GabDewraj/library-api/pkgs/api/middleware"
	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

type LibraryRouterParams struct {
	fx.In
	Mux        *chi.Mux
	Middleware middleware.Service
	Handler    handlers.BooksHandler
}

func NewBooksRouter(params LibraryRouterParams) error {
	params.Mux.Use(params.Middleware.CORS)
	// Routes
	params.Mux.Post("/books", params.Handler.CreateBook)
	params.Mux.Get("/books", params.Handler.GetBooks)
	params.Mux.Get("/books/{book_id}", params.Handler.GetBookByID)
	params.Mux.Put("/books/{book_id}", params.Handler.UpdateBook)
	params.Mux.Delete("/books/{book_id}", params.Handler.ArchiveBook)
	return nil
}
