package apps

import (
	"context"
	"os"

	"github.com/GabDewraj/library-api/cmd/config"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/cache"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/repo"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type LibraryParams struct {
	fx.In
	Cfg      *config.Config
	Router   *chi.Mux
	Logger   *logrus.Logger
	DB       *sqlx.DB
	Redis    *redis.Client
	CTX      context.Context
	Shutdown chan os.Signal
}

func LibraryApp(p LibraryParams) {
	// Create the application
	app := fx.New(
		fx.Supply(
			p.Router,
			p.DB,
			p.Cfg,
			p.Redis,
		),
		fx.Provide(
			cache.NewRedisCache,
			repo.NewlibraryDB,
		),
		// fx.Invoke(router.NewLibrarysRouter),
		// Initiliase all separate server applications
	)

	// Each fx child has its own dependency cycle with its own context
	// We simply need to shut down the whole server if an application cannot startup and log the error
	logrus.Infoln("Librarys application is running...")
	if err := app.Start(p.CTX); err != nil {
		logrus.Errorf("Librarys application is shutting down with ERR: %v", err)
		logrus.Error(err)
		os.Exit(1)
		return
	}
	// Wait for the shutdown signal
	go func() {
		<-p.Shutdown
		logger := logrus.StandardLogger()
		logger.Info("Received shutdown signal. Shutting down gracefully...")

		// Stop the application
		if err := app.Stop(p.CTX); err != nil {
			logger.Error("Error stopping the application:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
}
