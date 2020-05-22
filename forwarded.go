package serve

import (
	"crypto/tls"
	"net"
	"net/http"
)

// TODO: Support multiple values and "Forwarded" header?
//  => https://github.com/gorilla/handlers/blob/master/proxy_headers.go

// Forwarded is a middleware that will parse "X-Forwarded-X" headers and mutate
// the request to reflect the conditions described by the headers.
//
// Note: This technique should only be applied to apps that are behind a load
// balancer that will *always* set the headers. Otherwise an attacker may be
// able to provide false information and circumvent security limitations.
func Forwarded() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get headers
			ip := r.Header.Get("X-Forwarded-For")
			port := r.Header.Get("X-Forwarded-Port")
			proto := r.Header.Get("X-Forwarded-Proto")

			// rewrite remote addr
			if ip != "" && port != "" {
				r.RemoteAddr = net.JoinHostPort(ip, port)
			}

			// fake tls if scheme is https
			if r.TLS == nil && proto == "https" {
				// set url scheme
				r.URL.Scheme = "https"

				// set fake tls state
				r.TLS = &tls.ConnectionState{
					Version:           tls.VersionTLS13,
					HandshakeComplete: true,
					CipherSuite:       tls.TLS_AES_256_GCM_SHA384,
				}
			}

			// call next
			next.ServeHTTP(w, r)
		})
	}
}
