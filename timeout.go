package serve

import (
	"context"
	"net/http"
	"time"
)

// Timeout returns a middleware that ensures a request timeout.
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// add timeout to context
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// call next
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
