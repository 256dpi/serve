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

// CORSDefault returns a default cors policy for basic APIs.
func CORSDefault(headers ...string) CORSPolicy {
	return CORSPolicy{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowedHeaders: append([]string{
			"Accept", "Authorization", "Content-Type", "Origin",
			"X-Requested-With",
		}, headers...),
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}
}
