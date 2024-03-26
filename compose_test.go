package serve

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompose(t *testing.T) {
	assert.PanicsWithValue(t, `compose: expected chain to have at least two items`, func() {
		Compose()
	})

	assert.PanicsWithValue(t, `compose: expected last chain item to be a "http.Handler"`, func() {
		Compose(nil, nil)
	})

	assert.PanicsWithValue(t, `compose: expected intermediary chain item to be a "func(http.handler) http.Handler"`, func() {
		Compose(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	})

	handler := Compose(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("1"))
				next.ServeHTTP(w, r)
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("2"))
				next.ServeHTTP(w, r)
			})
		},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("H"))
		}),
	)

	r := Record(nil, handler, "GET", "/foo", nil, "")
	assert.Equal(t, "12H", r.Body.String())
}
