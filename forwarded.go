package serve

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// TODO: Support new "Forwarded" header?
//  => https://github.com/gorilla/handlers/blob/master/proxy_headers.go

// Forwarded is a middleware that will parse the selected "X-Forwarded-X" headers
// and mutate the request to reflect the conditions described by the headers.
//
// Note: This technique should only be applied to apps that are behind a load
// balancer that will *always* set the selected headers. Otherwise an attacker
// may be able to provide false information and circumvent security limitations.
func Forwarded(useIP, usePort, useProto, fakeTLS bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get ip, port and protocol
			ip, port, _ := net.SplitHostPort(r.RemoteAddr)
			protocol := r.URL.Scheme

			// get forwarded ip
			if useIP {
				forwardedIP := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
				if net.ParseIP(forwardedIP) != nil {
					ip = forwardedIP
				}
			}

			// get forwarded port
			if usePort {
				forwardedPort := r.Header.Get("X-Forwarded-Port")
				if n, _ := strconv.Atoi(forwardedPort); n > 0 {
					port = forwardedPort
				}
			}

			// get forwarded protocol
			if useProto {
				forwardedProtocol := r.Header.Get("X-Forwarded-Proto")
				if forwardedProtocol == "https" {
					protocol = "https"
				}
			}

			// rewrite remote addr if changed
			remote := net.JoinHostPort(ip, port)
			if r.RemoteAddr != remote {
				r.RemoteAddr = remote
			}

			// fake tls if scheme is https
			if fakeTLS && r.TLS == nil && protocol == "https" {
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
