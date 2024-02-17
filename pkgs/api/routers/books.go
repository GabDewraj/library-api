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
