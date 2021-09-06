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

// CORSDefault returns a default cors policy for basic APIs. Set origin to "*"
// to allow request from any origin.
func CORSDefault(origin string, headers ...string) CORSPolicy {
	return CORSPolicy{
		AllowedOrigins: []string{origin},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowedHeaders: append([]string{
			"Accept", "Authorization", "Cache-Control", "Content-Disposition",
			"Content-Type", "Origin", "X-Requested-With",
		}, headers...),
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}
}
