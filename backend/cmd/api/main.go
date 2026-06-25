package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/a3dif3x/yat/backend/internal/config"
	"github.com/a3dif3x/yat/backend/internal/middleware"
)

func main() {

	config, err := config.Load()
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("config load failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel}))

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthCheckHandler)

	handler := middleware.Chain(
		middleware.RequestID(),
		middleware.Recovery(logger),
		middleware.RequestLogging(logger),
	)(mux)

	address := ":" + config.Port

	logger.Info("started server", slog.String("address", address))
	if err := http.ListenAndServe(address, handler); err != nil {
		logger.Error("server stopped", slog.Any("error", err))
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
