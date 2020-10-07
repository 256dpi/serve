package serve

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	handler := Compose(
		BasicAuth("foo", "bar", "Test"),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Protected"))
		}),
	)

	r := Record(handler, "GET", "/foo", nil, "")
	assert.Equal(t, "", r.Body.String())
	assert.Equal(t, http.Header{
		"Www-Authenticate": []string{`Basic realm="Test"`},
	}, r.Header())

	r = Record(handler, "GET", "/foo", map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("foo:foo")),
	}, "")
	assert.Equal(t, "", r.Body.String())
	assert.Equal(t, http.Header{
		"Www-Authenticate": []string{`Basic realm="Test"`},
	}, r.Header())

	r = Record(handler, "GET", "/foo", map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar")),
	}, "")
	assert.Equal(t, "Protected", r.Body.String())
	assert.Equal(t, http.Header{
		"Content-Type": []string{"text/plain; charset=utf-8"},
	}, r.Header())
}
