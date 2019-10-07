package serve

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecurity(t *testing.T) {
	sec := Compose(
		Security(false, false, 7*24*time.Hour),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello"))
		}),
	)

	r := Record(sec, "GET", "http://example.com", nil, "")
	assert.Equal(t, http.StatusMovedPermanently, r.Code)
	assert.Equal(t, "<a href=\"https://example.com\">Moved Permanently</a>.\n\n", r.Body.String())
	assert.Equal(t, http.Header{
		"Content-Type": []string{"text/html; charset=utf-8"},
		"Location":     []string{"https://example.com"},
	}, r.Header())

	r = Record(sec, "GET", "https://example.com", nil, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "Hello", r.Body.String())
	assert.Equal(t, http.Header{
		"Content-Security-Policy":   []string{"default-src 'none'"},
		"Content-Type":              []string{"text/plain; charset=utf-8"},
		"Referrer-Policy":           []string{"no-referrer-when-downgrade"},
		"Strict-Transport-Security": []string{"max-age=604800"},
		"X-Content-Type-Options":    []string{"nosniff"},
		"X-Frame-Options":           []string{"DENY"},
		"X-Xss-Protection":          []string{"1; mode=block"},
	}, r.Header())
}

func TestSecurityAllowInsecure(t *testing.T) {
	sec := Compose(
		Security(true, false, 7*24*time.Hour),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello"))
		}),
	)

	r := Record(sec, "GET", "http://example.com", nil, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "Hello", r.Body.String())
	assert.Equal(t, http.Header{
		"Content-Security-Policy":   []string{"default-src 'none'"},
		"Content-Type":              []string{"text/plain; charset=utf-8"},
		"Referrer-Policy":           []string{"no-referrer-when-downgrade"},
		"Strict-Transport-Security": []string{"max-age=604800"},
		"X-Content-Type-Options":    []string{"nosniff"},
		"X-Frame-Options":           []string{"DENY"},
		"X-Xss-Protection":          []string{"1; mode=block"},
	}, r.Header())
}

func TestSecurityNoFrontend(t *testing.T) {
	sec := Compose(
		Security(false, true, 7*24*time.Hour),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello"))
		}),
	)

	r := Record(sec, "GET", "https://example.com", nil, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "Hello", r.Body.String())
	assert.Equal(t, http.Header{
		"Content-Security-Policy":   []string{"default-src 'none'"},
		"Content-Type":              []string{"text/plain; charset=utf-8"},
		"Referrer-Policy":           []string{"origin-when-cross-origin, strict-origin-when-cross-origin"},
		"Strict-Transport-Security": []string{"max-age=604800"},
		"X-Content-Type-Options":    []string{"nosniff"},
		"X-Frame-Options":           []string{"DENY"},
		"X-Xss-Protection":          []string{"1; mode=block"},
	}, r.Header())
}
