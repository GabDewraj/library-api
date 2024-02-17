package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GabDewraj/library-api/cmd/apps"
	"github.com/GabDewraj/library-api/cmd/config"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	rootCmd cobra.Command
)

func main() {
	// Test migrations
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "test-migration",
			Short: "Creates a testDB within github actions",
			Long:  ``,
			Run: func(cmd *cobra.Command, args []string) {
				config.NewDBConnection(&config.Config{
					DB: config.DBConfig{
						Driver:                 "mysql",
						Host:                   "localhost",
						Port:                   "3306",
						Database:               "library_dev",
						Password:               "password",
						Username:               "root",
						MigrationDirectoryPath: "./cmd/server/config/migrations",
						ForceTLS:               false,
					},
				})
			},
		})

	// library Server
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run the library server",
			Long:  ``,
			Run: func(cmd *cobra.Command, args []string) {
				// Create a context to handle binary startup
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				// Handle OS signals to gracefully shutdown the application
				shutdownSignal := make(chan os.Signal, 1)
				signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
				app := fx.New(
					fx.Provide(
						func() (context.Context, chan os.Signal) {
							return ctx, shutdownSignal
						},
						logrus.StandardLogger,
						config.NewConfig,
						chi.NewRouter,
						config.NewDBConnection,
						config.NewRedisClient,
					),
					// Initialize all separate server applications
					fx.Invoke(apps.LibraryApp),
					// Run the router
					fx.Invoke(
						func(r *chi.Mux, cfg *config.Config, logger *logrus.Logger) {
							logger.Info("Server is running on port", cfg.ServerPort, "...")
							http.ListenAndServe(cfg.ServerPort, r)
						},
					),
				)

				// Pass child context to allow for
				go func(ctx context.Context) {
					// Wait for the shutdown signal
					<-shutdownSignal
					logger := logrus.StandardLogger()
					logger.Info("Received shutdown signal. Shutting down gracefully...")

					// Stop the application
					if err := app.Stop(ctx); err != nil {
						logger.Error("Error stopping the application:", err)
						os.Exit(1)
					}
					os.Exit(0)
				}(ctx)

				// Start the application
				if err := app.Start(ctx); err != nil {
					cancel()
					logrus.StandardLogger().Fatal("Error starting the application:", err)
				}
			},
		})
	rootCmd.Execute()
}
