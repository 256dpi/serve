package serve

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentSecurity(t *testing.T) {
	handler := Compose(
		ContentSecurity(ContentPolicy{
			"default-src": []string{"'self'", "https://example.com"},
			"style-src":   []string{"'self'", "'unsafe-inline'"},
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello"))
		}),
	)

	r := Record(nil, handler, "GET", "https://example.com", nil, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "Hello", r.Body.String())
	assert.Equal(t, http.Header{
		"Content-Security-Policy": []string{"default-src 'self' https://example.com; style-src 'self' 'unsafe-inline'"},
		"Content-Type":            []string{"text/plain; charset=utf-8"},
	}, r.Header())
}
