package serve

import (
	"net/http"
	"strings"
)

// AssetServer constructs an asset server handler that serves an asset
// directory on a specified path and serves the index file for not found paths
// which is needed to properly serve single page applications.
func AssetServer(prefix, directory string) http.Handler {
	// ensure prefix
	prefix = "/" + strings.Trim(prefix, "/")

	// create dir server
	dir := http.Dir(directory)

	// create file server
	fs := http.FileServer(dir)

	h := func(w http.ResponseWriter, r *http.Request) {
		// pre-check if file does exist
		f, err := dir.Open(r.URL.Path)
		if err != nil {
			r.URL.Path = "/"
		} else if f != nil {
			_ = f.Close()
		}

		// serve file
		fs.ServeHTTP(w, r)
	}

	return http.StripPrefix(prefix, http.HandlerFunc(h))
}
