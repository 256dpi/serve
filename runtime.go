package serve

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Runtime returns a middleware that ensures a minimum and maximum request
// runtime. If the minimum runtime is zero, only a maximum runtime is enforced.
func Runtime(min, max time.Duration) func(http.Handler) http.Handler {
	// handle zero minimum
	if min == 0 {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// add timeout to context
				ctx, cancel := context.WithTimeout(r.Context(), max)
				defer cancel()

				// call next
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// capture request context
			reqCtx := r.Context()

			// prepare wait group
			var wg sync.WaitGroup
			defer wg.Wait()

			// create a new context with timeout
			newCtx, cancel := context.WithTimeout(context.WithoutCancel(reqCtx), max)
			defer cancel()

			// cancel new context when request context is done, but only after
			// the specified minimum runtime has passed
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-newCtx.Done():
				case <-time.After(min):
					select {
					case <-newCtx.Done():
					case <-reqCtx.Done():
						cancel()
					}
				}
			}()

			// call next
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
