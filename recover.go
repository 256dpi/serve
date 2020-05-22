package serve

import (
	"errors"
	"fmt"
	"net/http"
)

// Recover is a middleware that recovers panics and forwards the error to the
// the provided reporter.
func Recover(reporter func(error)) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// remove all headers
					for key := range w.Header() {
						w.Header().Del(key)
					}

					// write header
					w.WriteHeader(http.StatusInternalServerError)

					// report error
					switch err := err.(type) {
					case error:
						reporter(err)
					case string:
						reporter(errors.New(err))
					default:
						reporter(fmt.Errorf("%v", err))
					}
				}
			}()

			// call next
			next.ServeHTTP(w, r)
		})
	}
}
