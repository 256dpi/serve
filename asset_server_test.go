package serve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetServer(t *testing.T) {
	as1 := AssetServer("/", ".test/assets/")

	r := testRequest(as1, "GET", "/", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = testRequest(as1, "GET", "/foo", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = testRequest(as1, "GET", "/foo/bar", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	as2 := AssetServer("/foo/", ".test/assets/")

	r = testRequest(as2, "GET", "/foo/", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = testRequest(as2, "GET", "/foo/foo", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())

	r = testRequest(as2, "GET", "/foo/bar", nil, "")
	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "<h1>Hello</h1>\n", r.Body.String())
}
