package middleware

import (
	"context"
	"crypto/rand"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return id
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := rand.Text()
			ctx := WithRequestID(r.Context(), id)
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
