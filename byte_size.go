package serve

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidByteSize is returned for invalid byte sizes.
var ErrInvalidByteSize = errors.New("serve: byte size must be like 4K, 20MiB or 5GB")

// MustByteSize will call ByteSize and panic on errors.
func MustByteSize(str string) int64 {
	// get byte size
	bytes, err := ByteSize(str)
	if err != nil {
		panic(err)
	}

	return bytes
}

// ByteSize parses human readable byte sizes (e.g. 4K, 20MiB or 5GB) and returns
// the amount of bytes they represent. ErrInvalidByteSize is returned if the
// specified byte size is invalid.
func ByteSize(str string) (int64, error) {
	// check length
	if len(str) < 2 {
		return 0, ErrInvalidByteSize
	}

	// get symbol index
	index := strings.IndexAny(str, "KMGT")
	if index <= 0 {
		return 0, ErrInvalidByteSize
	}

	// parse number
	num, err := strconv.ParseInt(str[:index], 10, 64)
	if err != nil {
		return 0, ErrInvalidByteSize
	}

	// calculate size
	switch str[index:] {
	case "K", "KiB":
		return num * 1000, nil
	case "KB":
		return num * 1024, nil
	case "M", "MiB":
		return num * 1000 * 1000, nil
	case "MB":
		return num * 1024 * 1024, nil
	case "G", "GiB":
		return num * 1000 * 1000 * 1000, nil
	case "GB":
		return num * 1024 * 1024 * 1024, nil
	case "T", "TiB":
		return num * 1000 * 1000 * 1000 * 1000, nil
	case "TB":
		return num * 1024 * 1024 * 1024 * 1024, nil
	}

	return 0, ErrInvalidByteSize
}
