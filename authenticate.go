package serve

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// Authenticate returns a middleware that enforces HTTP Basic Authentication.
func Authenticate(username, password, realm string) func(http.Handler) http.Handler {
	// hash username and password
	requiredUser := sha256.Sum256([]byte(username))
	requiredPass := sha256.Sum256([]byte(password))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get username and password
			username, password, ok := r.BasicAuth()

			// hash username and password
			givenUser := sha256.Sum256([]byte(username))
			givenPass := sha256.Sum256([]byte(password))

			// call next handler if ok
			if ok && subtle.ConstantTimeCompare(givenUser[:], requiredUser[:]) == 1 && subtle.ConstantTimeCompare(givenPass[:], requiredPass[:]) == 1 {
				next.ServeHTTP(w, r)
				return
			}

			// otherwise, require authentication
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
