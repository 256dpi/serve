package serve

import (
	"encoding/json"
	"net/http"
)

// RPCContext allows handlers to interact with the request.
type RPCContext struct {
	r *http.Request
	w http.ResponseWriter
}

// Request will return the original request.
func (c *RPCContext) Request() *http.Request {
	return c.r
}

// ResponseWriter will return the associated response writer.
func (c *RPCContext) ResponseWriter() http.ResponseWriter {
	return c.w
}

// Parse will parse the request to the specified value.
func (c *RPCContext) Parse(v interface{}) error {
	return json.NewDecoder(c.r.Body).Decode(&v)
}

// Handle will parse the request and yield control to the callback.
func (c *RPCContext) Handle(v interface{}, cb func() interface{}) interface{} {
	// parse request
	err := c.Parse(v)
	if err != nil {
		return err
	}

	return cb()
}

// RPCHandler wraps a handler to simplify handling request and responses.
func RPCHandler(handler func(*RPCContext) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// run handler
		res := handler(&RPCContext{r: r, w: w})
		if err, ok := res.(error); ok {
			// check error
			rpcError, ok := err.(*RPCError)
			if !ok {
				rpcError = RPCInternalServerError("unknown error")
			}

			// set status
			if http.StatusText(rpcError.Status) == "" {
				rpcError.Status = http.StatusInternalServerError
			}

			// write header
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(rpcError.Status)

			// write error
			_ = json.NewEncoder(w).Encode(rpcError)

			return
		}

		// handle nil response
		if res == nil {
			res = RPCData(nil)
		}

		// write header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// write response
		_ = json.NewEncoder(w).Encode(res)
	}
}
