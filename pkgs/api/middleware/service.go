package middleware

import (
	"net/http"

	"github.com/GabDewraj/library-api/pkgs/infrastructure/cache"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Cache cache.Service
}

type service struct {
	Cache cache.Service
}

type Service interface {
	CORS(next http.Handler) http.Handler
	RateLimiter(next http.Handler) http.Handler
}

func NewMiddlwareStack(p Params) Service {
	return &service{
		Cache: p.Cache,
	}
}
