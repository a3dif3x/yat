package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/a3dif3x/yat/backend/internal/middleware"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthCheckHandler)

	handler := middleware.Chain(
		middleware.RequestID(),
		middleware.Recovery(logger),
		middleware.RequestLogging(logger),
	)(mux)

	logger.Info("started server", slog.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Error("server stopped", slog.Any("error", err))
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
