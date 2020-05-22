package serve

import (
	"net"
	"net/url"
	"strings"
)

// Hostname will return the hostname from the provided host string. This method
// should be used instead of net.SplitHostPort when attempting to clean the
// http/Request.Host attribute.
func Hostname(host string) string {
	return (&url.URL{Host: host}).Hostname()
}

// IP will return just the IP part from an address of the form ip[:port].
func IP(addr string) string {
	// check colon
	if !strings.ContainsRune(addr, ':') {
		return addr
	}

	// split addr
	ip, _, _ := net.SplitHostPort(addr)

	return ip
}
