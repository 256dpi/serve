package serve

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForwarded(t *testing.T) {
	sec := Compose(
		Forwarded(),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			str := fmt.Sprintf("%s: URL: %s, TLS: %v", r.RemoteAddr, r.URL.String(), r.TLS != nil)
			_, _ = w.Write([]byte(str))
		}),
	)

	r := Record(sec, "GET", "http://example.com", nil, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "192.0.2.1:1234: URL: http://example.com, TLS: false", r.Body.String())

	r = Record(sec, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For":  "1.2.3.4.",
		"X-Forwarded-Port": "4321",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "1.2.3.4.:4321: URL: http://example.com, TLS: false", r.Body.String())

	r = Record(sec, "GET", "http://example.com", map[string]string{
		"X-Forwarded-Proto": "https",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "192.0.2.1:1234: URL: https://example.com, TLS: true", r.Body.String())
}
