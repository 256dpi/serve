package serve

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	handler := Compose(
		Timeout(time.Millisecond),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
			w.WriteHeader(http.StatusTooManyRequests)
		}),
	)

	res := Record(nil, handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusTooManyRequests, res.Code)
}
