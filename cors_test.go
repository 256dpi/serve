package serve

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	handler := Compose(
		CORS(CORSPolicy{
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"Content-Type"},
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	res := Record(nil, handler, "GET", "/", map[string]string{
		"Origin": "example.com",
	}, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, http.Header{
		"Access-Control-Allow-Origin": []string{"*"},
		"Vary":                        []string{"Origin"},
	}, res.Header())

	res = Record(nil, handler, "OPTIONS", "/", map[string]string{
		"Origin": "example.com",
	}, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, http.Header{
		"Access-Control-Allow-Origin": []string{"*"},
		"Vary":                        []string{"Origin"},
	}, res.Header())

	res = Record(nil, handler, "OPTIONS", "/", map[string]string{
		"Origin":                        "example.com",
		"Access-Control-Request-Method": "POST",
	}, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, http.Header{
		"Access-Control-Allow-Origin":  []string{"*"},
		"Access-Control-Allow-Methods": []string{"POST"},
		"Vary":                         []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
	}, res.Header())
}

func TestCORSDefault(t *testing.T) {
	handler := Compose(
		CORS(CORSDefault("example.com")),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	res := Record(nil, handler, "GET", "/", map[string]string{
		"Origin": "example.com",
	}, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, http.Header{
		"Access-Control-Allow-Origin":      {"example.com"},
		"Access-Control-Allow-Credentials": {"true"},
		"Access-Control-Expose-Headers":    {"Content-Type"},
		"Vary":                             {"Origin"},
	}, res.Header())

	res = Record(nil, handler, "GET", "/", map[string]string{
		"Origin": "bar.com",
	}, "")
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, http.Header{
		"Vary": {"Origin"},
	}, res.Header())
}
