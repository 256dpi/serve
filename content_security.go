package serve

import (
	"net/http"
	"sort"
	"strings"
)

// TODO: Merge-able content security policies?

// ContentPolicy for defining content security.
type ContentPolicy map[string][]string

// String will encode the policy as a string.
func (p ContentPolicy) String() string {
	// collect segments
	segments := make([]string, 0, len(p))
	for directive, sources := range p {
		segments = append(segments, directive+" "+strings.Join(sources, " "))
	}

	// sort segments
	sort.Strings(segments)

	return strings.Join(segments, "; ")
}

// ContentSecurity returns a middleware for enforcing content security.
func ContentSecurity(policy ContentPolicy) func(http.Handler) http.Handler {
	// precompile policy
	compiledPolicy := policy.String()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// set header
			w.Header().Set("Content-Security-Policy", compiledPolicy)

			// call next
			next.ServeHTTP(w, r)
		})
	}
}
