package serve

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	errs := make(chan error, 1)

	sec := Compose(
		Recover(func(err error) {
			errs <- err
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Foo", "Bar")

			if r.Method == "POST" {
				panic(fmt.Errorf("foo"))
			} else if r.Method == "DELETE" {
				panic(map[string]int{
					"n": 42,
				})
			}

			_, _ = w.Write([]byte("Hello"))
		}),
	)

	r := Record(sec, "GET", "http://example.com", nil, "")
	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "Hello", r.Body.String())
	assert.Equal(t, http.Header{
		"Content-Type": {"text/plain; charset=utf-8"},
		"Foo":          {"Bar"},
	}, r.Header())

	r = Record(sec, "POST", "http://example.com", nil, "")
	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "", r.Body.String())
	assert.Equal(t, http.Header{}, r.Header())
	err := <-errs
	assert.Error(t, err)
	assert.Equal(t, "foo", err.Error())

	r = Record(sec, "DELETE", "http://example.com", nil, "")
	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "", r.Body.String())
	assert.Equal(t, http.Header{}, r.Header())
	err = <-errs
	assert.Error(t, err)
	assert.Equal(t, "map[n:42]", err.Error())
}
