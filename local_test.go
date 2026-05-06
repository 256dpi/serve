package serve

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocal(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK!"))
	})

	client := http.Client{
		Transport: Local(handler),
	}

	res, err := client.Get("http://0.0.0.0/foo")
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK!", string(data))
}
