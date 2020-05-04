package serve

import "net/http"

// Throttle returns a middleware that limits concurrent requests. The middleware
// will block the request and wait for a token. If the context is cancelled
// beforehand, "Too Many Requests" is returned.
func Throttle(concurrency int) func(http.Handler) http.Handler {
	// create bucket
	bucket := make(chan struct{}, concurrency)

	// fill bucket
	for i := 0; i < concurrency; i++ {
		bucket <- struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get token
			select {
			case <-bucket:
			case <-r.Context().Done():
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			// ensure token is added back
			defer func() {
				select {
				case bucket <- struct{}{}:
				default:
				}
			}()

			// call next
			next.ServeHTTP(w, r)
		})
	}
}
