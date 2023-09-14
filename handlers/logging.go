package handlers

import (
	"net/http"
	"time"
)

func (h *HandlerHelper) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the provided handler function
		next.ServeHTTP(w, r)

		// Logging the request details
		h.l.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Dur("duration", time.Since(start)).
			Msg("request handled")
	})
}
