package middleware

import (
	"net/http"
	"time"

	"github.com/GabDewraj/library-api/cmd/config"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/cache"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Cache  cache.Service
	Config *config.Config
}

type service struct {
	Cache                cache.Service
	RateWindow           time.Duration
	MaxRequestsPerWindow int
}

type Service interface {
	CORS(next http.Handler) http.Handler
	RateLimiter(next http.Handler) http.Handler
	CustomLogger(next http.Handler) http.Handler
}

func NewMiddlwareStack(p Params) Service {
	middleware := p.Config.MiddlewareConfig
	return &service{
		Cache:                p.Cache,
		MaxRequestsPerWindow: middleware.MaxRequests,
		RateWindow:           time.Duration(middleware.RateWindow) * time.Minute,
	}
}
