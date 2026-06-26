package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a3dif3x/yat/backend/internal/config"
	"github.com/a3dif3x/yat/backend/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: config.LogLevel}))

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthCheckHandler)

	handler := middleware.Chain(
		middleware.RequestID(),
		middleware.Recovery(logger),
		middleware.RequestLogging(logger),
	)(mux)

	srv := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverError := make(chan error, 1)

	go func() {
		logger.Info("started server", slog.String("address", srv.Addr))
		serverError <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	logger.Info("server stopped cleanly")
	return nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func readyCheckHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("NOT READY"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("READY"))
	}
}
