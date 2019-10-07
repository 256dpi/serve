package serve

import (
	"net/http"
	"strconv"
	"time"
)

// Security will return a middleware that enforce various standard web security
// techniques.
func Security(allowInsecure, noFrontend bool, stsMaxAge time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// call next handle if request is secure
			if Secure(w, r, allowInsecure, noFrontend, stsMaxAge) {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// Secure will enforce various common web security policies. It return whether
// the request is safe to be further processed. Subsequent handlers should
// update the headers with a more applicable content security policy.
func Secure(w http.ResponseWriter, r *http.Request, allowInsecure, noFrontend bool, stsMaxAge time.Duration) bool {
	// redirect insecure request if not allowed and not secure
	if !allowInsecure && r.TLS == nil {
		url := *r.URL
		url.Host = r.Host
		url.Scheme = "https"
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		return false
	}

	// set basic security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Xss-Protection", "1; mode=block")

	// set referrer policy
	if noFrontend {
		// only send origin for cross origin requests
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin, strict-origin-when-cross-origin")
	} else {
		// no referrer when downgraded or insecure
		w.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
	}

	// set strict transport max age if specified
	if stsMaxAge > 0 {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.Itoa(int(stsMaxAge/time.Second)))
	}

	// set default content security policy
	w.Header().Set("Content-Security-Policy", "default-src 'none'")

	return true
}
