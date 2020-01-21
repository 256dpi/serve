package serve

import (
	"errors"
	"io"
	"net/http"
)

// ErrBodyLimitExceeded is returned if a body is read beyond the set limit.
var ErrBodyLimitExceeded = errors.New("body limit exceeded")

// BodyLimiter wraps a io.ReadCloser and keeps a reference to the original.
type BodyLimiter struct {
	Length   int64
	Limit    int64
	Original io.ReadCloser
	Limited  io.ReadCloser
}

// LimitBody will limit reading from the body of the supplied request to the
// specified amount of bytes. Earlier calls to LimitBody will be overwritten
// which essentially allows callers to increase the limit from a default limit
// later in the request processing.
func LimitBody(w http.ResponseWriter, r *http.Request, limit int64) {
	// recover original body from existing limiter
	if bl, ok := r.Body.(*BodyLimiter); ok {
		r.Body = bl.Original
	}

	// set limited body
	r.Body = &BodyLimiter{
		Length:   r.ContentLength,
		Limit:    limit,
		Original: r.Body,
		Limited:  http.MaxBytesReader(w, r.Body, limit),
	}
}

// Read will read from the underlying io.Reader.
func (l *BodyLimiter) Read(p []byte) (int, error) {
	// immediately return error if length is beyond limit
	if l.Length >= 0 && l.Length > l.Limit {
		return 0, ErrBodyLimitExceeded
	}

	// read and rewrite error
	n, err := l.Limited.Read(p)
	if err != nil && err.Error() == "http: request body too large" {
		return n, ErrBodyLimitExceeded
	}

	return n, err
}

// Close will close the body.
func (l *BodyLimiter) Close() error {
	return l.Limited.Close()
}
