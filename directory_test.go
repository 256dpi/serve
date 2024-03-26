package serve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectory(t *testing.T) {
	handler := Directory("/", ".test/assets/")

	r := Record(nil, handler, "GET", "/", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = Record(nil, handler, "GET", "/foo", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = Record(nil, handler, "GET", "/foo/bar", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	handler = Directory("/foo/", ".test/assets/")

	r = Record(nil, handler, "GET", "/foo/", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = Record(nil, handler, "GET", "/foo/foo", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = Record(nil, handler, "GET", "/foo/bar", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())
}
