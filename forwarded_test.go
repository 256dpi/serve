package serve

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleCloud(t *testing.T) {
	Forwarded(GoogleCloud(false))
}

func TestForwarded(t *testing.T) {
	handler := Compose(
		Forwarded(false, false, false, false, 0),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			str := fmt.Sprintf("%s: URL: %s, TLS: %v", r.RemoteAddr, r.URL.String(), r.TLS != nil)
			_, _ = w.Write([]byte(str))
		}),
	)

	// ignored
	r := Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For":   "1.2.3.4",
		"X-Forwarded-Port":  "4321",
		"X-Forwarded-Proto": "https",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "192.0.2.1:1234: URL: http://example.com, TLS: false", r.Body.String())

	handler = Compose(
		Forwarded(true, true, true, true, 0),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			str := fmt.Sprintf("%s: URL: %s, TLS: %v", r.RemoteAddr, r.URL.String(), r.TLS != nil)
			_, _ = w.Write([]byte(str))
		}),
	)

	// missing
	r = Record(handler, "GET", "http://example.com", nil, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "192.0.2.1:1234: URL: http://example.com, TLS: false", r.Body.String())

	// correct
	r = Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For":   "1.2.3.4",
		"X-Forwarded-Port":  "4321",
		"X-Forwarded-Proto": "https",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "1.2.3.4:4321: URL: https://example.com, TLS: true", r.Body.String())

	// invalid
	r = Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For":   "foo",
		"X-Forwarded-Port":  "foo",
		"X-Forwarded-Proto": "foo",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "192.0.2.1:1234: URL: http://example.com, TLS: false", r.Body.String())

	// empty
	r = Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For":   "",
		"X-Forwarded-Port":  "",
		"X-Forwarded-Proto": "",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "192.0.2.1:1234: URL: http://example.com, TLS: false", r.Body.String())

	r = Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For": "2.3.4.5, 1.2.3.4",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "2.3.4.5:1234: URL: http://example.com, TLS: false", r.Body.String())

	handler = Compose(
		Forwarded(true, true, true, true, -2),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			str := fmt.Sprintf("%s: URL: %s, TLS: %v", r.RemoteAddr, r.URL.String(), r.TLS != nil)
			_, _ = w.Write([]byte(str))
		}),
	)

	r = Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For": "1.2.3.4",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "192.0.2.1:1234: URL: http://example.com, TLS: false", r.Body.String())

	r = Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For": "2.3.4.5, 1.2.3.4",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "2.3.4.5:1234: URL: http://example.com, TLS: false", r.Body.String())

	r = Record(handler, "GET", "http://example.com", map[string]string{
		"X-Forwarded-For": "3.4.5.6, 2.3.4.5, 1.2.3.4",
	}, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "2.3.4.5:1234: URL: http://example.com, TLS: false", r.Body.String())
}
