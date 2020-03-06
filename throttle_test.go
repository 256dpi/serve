package serve

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottle(t *testing.T) {
	handler := Compose(
		Timeout(5*time.Millisecond),
		Throttle(1),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
		}),
	)

	go func() {
		res := Record(handler, "GET", "/", nil, "")
		assert.Equal(t, http.StatusOK, res.Code)
	}()
	time.Sleep(time.Millisecond)

	res := Record(handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusTooManyRequests, res.Code)
}
