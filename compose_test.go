package serve

import (
	"net/http"
	"net/http/httptest"
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

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("H"))
	})

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("1"))
			next.ServeHTTP(w, r)
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("2"))
			next.ServeHTTP(w, r)
		})
	}

	e := Compose(m1, m2, h)

	r, err := http.NewRequest("GET", "/foo", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	e.ServeHTTP(w, r)
	assert.Equal(t, "12H", w.Body.String())
}
