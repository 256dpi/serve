package serve

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Record will make a request against the specified handler and record the result.
func Record(h http.Handler, method, url string, headers map[string]string, payload string) *httptest.ResponseRecorder {
	// prepare body
	var body io.Reader
	if payload != "" {
		body = strings.NewReader(payload)
	}

	// create request and recorder
	r := httptest.NewRequest(method, url, body)
	w := httptest.NewRecorder()

	// set headers
	for k, v := range headers {
		r.Header.Set(k, v)
	}

	// call handler
	h.ServeHTTP(w, r)

	return w
}
