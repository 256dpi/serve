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

// RPCHandler wraps a handler to simplify handling request and responses. The
// specified limit will be applied to the received request body. The handler is
// expected to return nil, JSON compatible responses, RPCError values or errors.
func RPCHandler(limit uint64, reporter func(error), handler func(*RPCContext) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check request method
		if r.Method != "POST" {
			RPCErrorWrite(w, RPCErrorFromStatus(http.StatusMethodNotAllowed, ""))
			return
		}

		// limit body
		LimitBody(w, r, limit)

		// run handler
		res := handler(&RPCContext{r: r, w: w})

		// handle errors
		if err, ok := res.(error); ok {
			// check error
			rpcError, ok := err.(*RPCError)
			if !ok {
				// report critical errors
				if reporter != nil {
					reporter(err)
				}

				// make rpc error
				rpcError = RPCErrorFromStatus(http.StatusInternalServerError, "")
			}

			// set status
			if http.StatusText(rpcError.Status) == "" {
				// report invalid errors
				if reporter != nil {
					reporter(rpcError)
				}

				// set fallback status
				rpcError.Status = http.StatusInternalServerError
			}

			// write error
			RPCErrorWrite(w, rpcError)

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
