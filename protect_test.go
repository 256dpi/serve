package serve

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProtect(t *testing.T) {
	handler := Compose(
		Protect(10, time.Second),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	for i := 0; i < 15; i++ {
		r := Record(nil, handler, "GET", "http://example.com", nil, "")
		if i <= 10 {
			assert.Equal(t, http.StatusOK, r.Code, i)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, r.Code, i)
		}
	}
}
