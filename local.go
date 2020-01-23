package serve

import (
	"net/http"
	"net/http/httptest"
)

type local struct {
	handler http.Handler
}

func (l *local) RoundTrip(req *http.Request) (*http.Response, error) {
	// prepare recorder
	rec := httptest.NewRecorder()

	// serve request
	l.handler.ServeHTTP(rec, req)

	return rec.Result(), nil
}

// Local returns a round tripper that uses the provided handler to serve the
// requests. It may be used with http.Client in unit tests.
func Local(handler http.Handler) http.RoundTripper {
	return &local{handler: handler}
}
