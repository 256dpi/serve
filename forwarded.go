package serve

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// GoogleCloud can be used with Forwarded to setup proper header parsing for
// traffic from Google Cloud load balancers.
func GoogleCloud(fakeTLS bool) (bool, bool, bool, bool, int) {
	return true, false, true, fakeTLS, -2
}

// Forwarded is a middleware that will parse the selected "X-Forwarded-X" headers
// and mutate the request to reflect the conditions described by the headers. As
// the "X-Forwarded-For" header may contain multiple values, the relative index
// of the client IP address must be specified.
//
// Note: This technique should only be applied to apps that are behind a load
// balancer that will *always* set/append the selected headers. Otherwise, an
// attacker may be able to provide false information and circumvent security
// limitations.
func Forwarded(useFor, usePort, useProto, fakeTLS bool, forIndex int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get ip, port and protocol
			ip, port, _ := net.SplitHostPort(r.RemoteAddr)
			protocol := r.URL.Scheme

			// get forwarded for
			if useFor {
				// get header
				forwardedFor := strings.Split(r.Header.Get("X-Forwarded-For"), ",")

				// compute index
				forwardedIndex := forIndex
				if forwardedIndex < 0 {
					forwardedIndex = len(forwardedFor) + forwardedIndex
				}

				// check bounds
				if forwardedIndex >= 0 && forwardedIndex < len(forwardedFor) {
					forwardedIP := strings.TrimSpace(forwardedFor[forwardedIndex])
					if net.ParseIP(forwardedIP) != nil {
						ip = forwardedIP
					}
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

			// update scheme if changed
			if r.URL.Scheme != protocol {
				r.URL.Scheme = protocol
			}

			// fake tls if scheme is https
			if fakeTLS && r.TLS == nil && protocol == "https" {
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
