package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

// @title           Library API
// @version         1.0
// @description     This is a sample library server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Gabriel Dewraj
// @contact.url    https://www.linkedin.com/in/gabriel-dewraj-8061681a2/
// @contact.email  gdewraj@gmail.com

// @host      localhost:8080

func main() {
	// Test migrations for gh actions
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "test-migration",
			Short: "Creates a testDB within github actions",
			Long:  ``,
			Run: func(cmd *cobra.Command, args []string) {
				// Create a context to handle binary startup
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				// Create test only config for safety, without reliance on infrastructure
				testConfig := &config.Config{
					DB: config.DBConfig{
						Driver:                 "mysql",
						Host:                   "localhost",
						Port:                   "3306",
						Database:               "library_dev",
						Password:               "password",
						Username:               "root",
						MigrationDirectoryPath: "./cmd/config/migrations",
						ForceTLS:               false,
					},
				}
				app := fx.New(
					// Supply the test config
					fx.Supply(testConfig),
					// Provide global server items to all applications
					fx.Provide(
						logrus.StandardLogger,
						config.NewDBConnection,
					),
					// Run necessary migrations
					fx.Invoke(config.PerformMigrations))
				// Start the server
				if err := app.Start(ctx); err != nil {
					cancel()
					logrus.StandardLogger().Fatal("Error running migrations", err)
				}
			},
		})

	// library Server
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run the library server",
			Long:  ``,
			Run: func(cmd *cobra.Command, args []string) {
				// Add mutex for shut down goroutines to avoid race conditions on shutdown
				mu := sync.Mutex{}
				// Create a context to handle binary startup
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				// Handle OS signals to gracefully shutdown the server
				shutdownSignal := make(chan os.Signal, 1)
				signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
				app := fx.New(
					// Provide global server items to all applications
					fx.Provide(
						func() (context.Context, chan os.Signal, *sync.Mutex) {
							return ctx, shutdownSignal, &mu
						},
						logrus.StandardLogger,
						config.NewConfig,
						chi.NewRouter,
						config.NewDBConnection,
						config.NewRedisClient,
					),
					// Run necessary migrations
					fx.Invoke(config.PerformMigrations),
					// Initialize all separate server applications
					fx.Invoke(apps.BooksApp),
					// Run the router
					fx.Invoke(
						func(r *chi.Mux, cfg *config.Config, logger *logrus.Logger) {
							logger.Info("Server is running on port", cfg.ServerPort, "...")
							http.ListenAndServe(cfg.ServerPort, r)
						},
					),
				)

				// Pass child context to allow for shutdown
				go func(ctx context.Context, mu *sync.Mutex) {
					mu.Lock()
					// Wait for the shutdown signal
					<-shutdownSignal
					logger := logrus.StandardLogger()
					logger.Info("Received shutdown signal. Shutting down gracefully...")

					// Stop the server
					if err := app.Stop(ctx); err != nil {
						logger.Error("Error stopping the server:", err)
						os.Exit(1)
					}
					mu.Unlock() // Avoid dead lock before the shutdown
					os.Exit(0)

				}(ctx, &mu)

				// Start the server
				if err := app.Start(ctx); err != nil {
					cancel()
					logrus.StandardLogger().Fatal("Error starting the server:", err)
				}
			},
		})
	rootCmd.Execute()
}
