package serve

import (
	"errors"
	"io"
	"net/http"
	"strconv"
)

// ErrBodyLimitExceeded is returned if a body is read beyond the set limit.
var ErrBodyLimitExceeded = errors.New("body limit exceeded")

// BodyLimiter wraps a io.ReadCloser and keeps a reference to the original.
type BodyLimiter struct {
	Limit    uint64
	Original io.ReadCloser
	Limited  io.ReadCloser
}

// LimitBody will limit reading from the body of the supplied request to the
// specified amount of bytes. Earlier calls to LimitBody will be overwritten
// which essentially allows callers to increase the limit from a default limit
// later in the request processing.
func LimitBody(w http.ResponseWriter, r *http.Request, limit uint64) {
	// get original body from existing limiter
	if bl, ok := r.Body.(*BodyLimiter); ok {
		r.Body = bl.Original
	}

	// set new limiter
	r.Body = &BodyLimiter{
		Original: r.Body,
		Limit:    limit,
		Limited:  http.MaxBytesReader(w, r.Body, int64(limit)),
	}
}

// Read will read from the underlying io.Reader.
func (l *BodyLimiter) Read(p []byte) (int, error) {
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

// DataSize parses human readable data sizes (e.g. 4K, 20M or 5G) and returns
// the amount of bytes they represent.
func DataSize(str string) uint64 {
	const msg = "fire: data size must be like 4K, 20M or 5G"

	// check length
	if len(str) < 2 {
		panic(msg)
	}

	// get symbol
	sym := string(str[len(str)-1])

	// parse number
	num, err := strconv.ParseUint(str[:len(str)-1], 10, 64)
	if err != nil {
		panic(msg)
	}

	// calculate size
	switch sym {
	case "K":
		return num * 1000
	case "M":
		return num * 1000 * 1000
	case "G":
		return num * 1000 * 1000 * 1000
	}

	panic(msg)
}
