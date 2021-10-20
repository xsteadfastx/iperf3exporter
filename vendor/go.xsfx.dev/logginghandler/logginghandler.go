// Package logginghandler is a simple, zerolog based, request logging http middleware.
// It also sets `X-Request-ID` in the request and response headers.
package logginghandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GetUUID gets the requests UUID from a request.
func GetUUID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}

// Logger returns a logger with the UUID set.
func Logger(r *http.Request) zerolog.Logger {
	logger := log.With().Str("uuid", GetUUID(r)).Logger()

	return logger
}

// Handler is the http middleware handler.
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New().String()
		r.Header.Set("X-Request-ID", uuid)
		logger := Logger(r)
		logger.Info().
			Str("method", r.Method).
			Str("user-agent", r.UserAgent()).
			Str("proto", r.Proto).
			Str("referer", r.Referer()).
			Str("request-url", r.URL.String()).
			Str("remote", r.RemoteAddr).
			Msg("")

		w.Header().Set("X-Request-ID", uuid)
		next.ServeHTTP(w, r)
	})
}
