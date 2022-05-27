package serve

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// ForwardedConfig defines handling of "X-Forwarded-X" headers.
type ForwardedConfig struct {
	UseFor   bool
	UsePort  bool
	UseProto bool
	FakeTLS  bool
	ForIndex int
	Debug    bool
}

// GoogleCloud can be used with Forwarded to setup proper header parsing for
// traffic from Google Cloud load balancers.
func GoogleCloud(fakeTLS bool) ForwardedConfig {
	return ForwardedConfig{
		UseFor:   true,
		UsePort:  false,
		UseProto: true,
		FakeTLS:  fakeTLS,
		ForIndex: -2,
	}
}

// ParseForwardedConfig will parse a forwarded config from the specified string
// and return it. This function can be used to infer a configuration on runtime
// from an environment variable or configuration file. The following comma
// seperated list of keywords ist supported: "use-for", "use-port", "use-proto",
// "fake-tls" and "for-index=1".
func ParseForwardedConfig(str string) ForwardedConfig {
	// parse keywords
	keywords := map[string]string{}
	for _, kw := range strings.Split(str, ",") {
		kw = strings.TrimSpace(kw)
		kv := strings.Split(kw, "=")
		if len(kv) == 1 {
			keywords[kv[0]] = ""
		} else if len(kv) == 2 {
			keywords[kv[0]] = kv[1]
		}
	}

	// get keywords
	_, useFor := keywords["use-for"]
	_, usePort := keywords["use-port"]
	_, useProto := keywords["use-proto"]
	_, fakeTLS := keywords["fake-tls"]
	forIndex, _ := strconv.Atoi(keywords["for-index"])
	_, debug := keywords["debug"]

	return ForwardedConfig{
		UseFor:   useFor,
		UsePort:  usePort,
		UseProto: useProto,
		FakeTLS:  fakeTLS,
		ForIndex: forIndex,
		Debug:    debug,
	}
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
func Forwarded(config ForwardedConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get ip, port and protocol
			ip, port, _ := net.SplitHostPort(r.RemoteAddr)
			protocol := r.URL.Scheme

			// debug
			if config.Debug {
				fmt.Printf("serve: forwarded headers: for=%q, port=%q, proto=%q\n",
					r.Header.Get("X-Forwarded-For"),
					r.Header.Get("X-Forwarded-Port"),
					r.Header.Get("X-Forwarded-Proto"),
				)
			}

			// get forwarded for
			if config.UseFor {
				// get header
				forwardedFor := strings.Split(r.Header.Get("X-Forwarded-For"), ",")

				// compute index
				forwardedIndex := config.ForIndex
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
			if config.UsePort {
				forwardedPort := r.Header.Get("X-Forwarded-Port")
				if n, _ := strconv.Atoi(forwardedPort); n > 0 {
					port = forwardedPort
				}
			}

			// get forwarded protocol
			if config.UseProto {
				forwardedProtocol := r.Header.Get("X-Forwarded-Proto")
				if forwardedProtocol == "https" {
					protocol = "https"
				}
			}

			// rewrite remote addr if changed
			remote := net.JoinHostPort(ip, port)
			if r.RemoteAddr != remote {
				if config.Debug {
					fmt.Printf("==> changing remote addr from %q to %q\n", r.RemoteAddr, remote)
				}
				r.RemoteAddr = remote
			}

			// update scheme if changed
			if r.URL.Scheme != protocol {
				if config.Debug {
					fmt.Printf("==> changing url scheme from %q to %q\n", r.URL.Scheme, protocol)
				}
				r.URL.Scheme = protocol
			}

			// fake tls if scheme is https
			if config.FakeTLS && r.TLS == nil && protocol == "https" {
				if config.Debug {
					fmt.Println("=> faking TLS connection")
				}
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
