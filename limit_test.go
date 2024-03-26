package serve

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ io.ReadCloser = &BodyLimiter{}

func TestLimit(t *testing.T) {
	handler := Compose(
		Limit(10),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := ioutil.ReadAll(r.Body)
			if err == ErrBodyLimitExceeded {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
			}
		}),
	)

	res := Record(nil, handler, "GET", "/", nil, "Hello!")
	assert.Equal(t, http.StatusOK, res.Code)

	res = Record(nil, handler, "GET", "/", nil, "Hello World!")
	assert.Equal(t, http.StatusRequestEntityTooLarge, res.Code)
}

func TestLimitBodyExtend(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.org", strings.NewReader("Hello World!"))

	orig := r.Body

	LimitBody(nil, r, 2)
	assert.Equal(t, orig, r.Body.(*BodyLimiter).Original)

	LimitBody(nil, r, 16)
	assert.Equal(t, orig, r.Body.(*BodyLimiter).Original)

	bytes, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", string(bytes))
}

func TestLimitBodyBeforeReading(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.org", strings.NewReader("Hello World!"))

	LimitBody(nil, r, 5)

	bytes, err := ioutil.ReadAll(r.Body)
	assert.Error(t, err)
	assert.Equal(t, "", string(bytes))
	assert.Equal(t, err, ErrBodyLimitExceeded)
}

func TestLimitBodyWhileReading(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.org", strings.NewReader("Hello World!"))
	r.ContentLength = -1

	LimitBody(nil, r, 5)

	bytes, err := ioutil.ReadAll(r.Body)
	assert.Error(t, err)
	assert.Equal(t, "Hello", string(bytes))
	assert.Equal(t, err, ErrBodyLimitExceeded)
}
