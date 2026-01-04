package webserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle, cfg Config, router *Router, logger *slog.Logger) *http.Server {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server", "port", cfg.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
