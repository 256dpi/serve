package serve

import (
	"fmt"
	"net/http"
	"strings"
)

// RPCError objects provide additional information about problems encountered
// while performing an RPC operation.
type RPCError struct {
	// A unique identifier for this particular occurrence of the problem.
	ID string `json:"id,omitempty"`

	// An URL that leads to further details about this problem.
	Link string `json:"link,omitempty"`

	// The HTTP status code applicable to this problem.
	Status int `json:"status,string,omitempty"`

	// An application-specific error code.
	Code string `json:"code,omitempty"`

	// A short, human-readable summary of the problem.
	Title string `json:"title,omitempty"`

	// A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// A JSON pointer specifying the source of the error.
	//
	// See https://tools.ietf.org/html/rfc6901.
	Source string `json:"source,omitempty"`

	// Non-standard meta-information about the error.
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// RPCError returns a string representation of the error for logging purposes.
func (e *RPCError) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

// RPCErrorFromStatus will return an error that has been derived from the passed
// status code.
//
// Note: If the passed status code is not a valid HTTP status code, an Internal
// Server RPCError status code will be used instead.
func RPCErrorFromStatus(status int, detail string) *RPCError {
	// get text
	str := strings.ToLower(http.StatusText(status))

	// check text
	if str == "" {
		status = http.StatusInternalServerError
		str = strings.ToLower(http.StatusText(http.StatusInternalServerError))
	}

	return &RPCError{
		Status: status,
		Title:  str,
		Detail: detail,
	}
}

// RPCNotFound returns a new not found error.
func RPCNotFound(detail string) *RPCError {
	return RPCErrorFromStatus(http.StatusNotFound, detail)
}

// RPCBadRequest returns a new bad request error with a source.
func RPCBadRequest(detail, source string) *RPCError {
	err := RPCErrorFromStatus(http.StatusBadRequest, detail)
	err.Source = source

	return err
}

// RPCInternalServerError returns na new internal server error.
func RPCInternalServerError(detail string) *RPCError {
	return RPCErrorFromStatus(http.StatusInternalServerError, detail)
}
