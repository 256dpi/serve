package serve

import (
	"net/http"
	"time"

	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

// Protect will return a middleware that will rate limit requests based on the
// remote IP address. It will allow up to the specified rate of requests per
// duration.
func Protect(rate int, duration time.Duration) func(http.Handler) http.Handler {
	// prepare store
	store, err := memstore.New(int(MustByteSize("100K")))
	if err != nil {
		panic(err)
	}

	// prepare rate limiter
	rateLimiter, err := throttled.NewGCRARateLimiter(store, throttled.RateQuota{
		MaxRate:  throttled.PerDuration(rate, duration),
		MaxBurst: rate,
	})
	if err != nil {
		panic(err)
	}

	// prepare handler
	handler := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy: &throttled.VaryBy{
			Custom: func(r *http.Request) string {
				return IP(r.RemoteAddr)
			},
		},
		DeniedHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}),
		Error: func(w http.ResponseWriter, r *http.Request, err error) {
			panic(err)
		},
	}

	return handler.RateLimit
}
