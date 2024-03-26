package serve

import (
	"context"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRuntime(t *testing.T) {
	routines := runtime.NumGoroutine()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusRequestTimeout)
			err := r.Context().Err()
			if err != nil {
				_, _ = w.Write([]byte(err.Error()))
			}
		case <-time.After(10 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		}
	})

	res := Record(nil, handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "", res.Body.String())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res = Record(ctx, handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusRequestTimeout, res.Code)
	assert.Equal(t, "context canceled", res.Body.String())

	/* with maximum only */

	handler = Runtime(0, time.Second)(handler)

	res = Record(nil, handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "", res.Body.String())

	ctx, cancel = context.WithCancel(context.Background())
	cancel()

	res = Record(ctx, handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusRequestTimeout, res.Code)
	assert.Equal(t, "context canceled", res.Body.String())

	/* with minimum and maximum */

	handler = Runtime(100*time.Millisecond, time.Second)(handler)

	res = Record(nil, handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "", res.Body.String())

	ctx, cancel = context.WithCancel(context.Background())
	cancel()

	res = Record(ctx, handler, "GET", "/", nil, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "", res.Body.String())

	assert.Equal(t, routines, runtime.NumGoroutine())
}
