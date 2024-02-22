package apps

import (
	"context"
	"os"
	"sync"

	"github.com/GabDewraj/library-api/cmd/config"
	"github.com/GabDewraj/library-api/pkgs/api/handlers"
	"github.com/GabDewraj/library-api/pkgs/api/middleware"
	"github.com/GabDewraj/library-api/pkgs/api/routers"
	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/cache/redcache"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/repo"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type BooksAppParams struct {
	fx.In
	Cfg      *config.Config
	Router   *chi.Mux
	Logger   *logrus.Logger
	DB       *sqlx.DB
	Redis    *redis.Client
	MU       *sync.Mutex
	CTX      context.Context
	Shutdown chan os.Signal
}

func BooksApp(p BooksAppParams) {
	// Create the application
	app := fx.New(
		fx.Supply(
			p.Router,
			p.DB,
			p.Cfg,
			p.Redis,
			p.MU,
		),
		fx.Provide(
			redcache.NewRedisCache,
			repo.NewBooksDB,
			books.NewService,
			middleware.NewMiddlwareStack,
			handlers.NewBooksHandler,
		),
		fx.Invoke(routers.NewBooksRouter),
	)

	logrus.Infoln("Books application is running...")
	if err := app.Start(p.CTX); err != nil {
		logrus.Errorf("Books application is shutting down with ERR: %v", err)
		logrus.Error(err)
		os.Exit(1)
		return
	}
	// Wait for the shutdown signal, using shared application to listen for cancel signal incase of error
	go func(ctx context.Context, mu *sync.Mutex) {
		mu.Lock()
		<-p.Shutdown
		logger := logrus.StandardLogger()
		logger.Info("Received shutdown signal. Shutting down gracefully...")

		// Stop the application
		if err := app.Stop(ctx); err != nil {
			logger.Error("Error stopping the application:", err)
			os.Exit(1)
		}
		mu.Unlock()
		os.Exit(0)
	}(p.CTX, p.MU)
}
