package serve

import "net/url"

// Hostname will return the hostname from the provided host string. This method
// should be used instead of net.SplitHostPort when attempting to clean the
// http/Request.Host attribute.
func Hostname(host string) string {
	return (&url.URL{Host: host}).Hostname()
}
