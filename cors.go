package serve

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSPolicy defines the CORS policy.
type CORSPolicy = cors.Options

// CORS returns a middleware for enforcing CORS.
func CORS(policy CORSPolicy) func(http.Handler) http.Handler {
	return cors.New(policy).Handler
}
