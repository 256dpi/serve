package serve

import (
	"io"
	"net/http"
	"strconv"
)

// BodyLimiter wraps a io.ReadCloser and keeps a reference to the original.
type BodyLimiter struct {
	io.ReadCloser
	Limit    uint64
	Original io.ReadCloser
}

// LimitBody will limit reading from the body of the supplied request to the
// specified amount of bytes. Earlier calls to LimitBody will be overwritten
// which essentially allows callers to increase the limit from a default limit.
func LimitBody(w http.ResponseWriter, r *http.Request, limit uint64) {
	// get original body from existing limiter
	if bl, ok := r.Body.(*BodyLimiter); ok {
		r.Body = bl.Original
	}

	// set new limiter
	r.Body = &BodyLimiter{
		Original:   r.Body,
		Limit:      limit,
		ReadCloser: http.MaxBytesReader(w, r.Body, int64(limit)),
	}
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
