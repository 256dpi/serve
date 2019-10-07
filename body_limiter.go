package serve

import (
	"io"
	"net/http"
)

// BodyLimiter wraps a io.ReadCloser and keeps a reference to the original.
type BodyLimiter struct {
	io.ReadCloser
	Limit    int64
	Original io.ReadCloser
}

// LimitBody will limit reading from the body of the supplied request to the
// specified amount of bytes. Earlier calls to LimitBody will be overwritten
// which essentially allows callers to increase the limit from a default limit.
func LimitBody(w http.ResponseWriter, r *http.Request, limit int64) {
	// get original body from existing limiter
	if bl, ok := r.Body.(*BodyLimiter); ok {
		r.Body = bl.Original
	}

	// set new limiter
	r.Body = &BodyLimiter{
		Original:   r.Body,
		Limit:      limit,
		ReadCloser: http.MaxBytesReader(w, r.Body, limit),
	}
}
