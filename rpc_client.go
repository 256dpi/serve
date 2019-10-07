package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RPCData represents arbitrary data.
type RPCData map[string]interface{}

// RPCClient is a reusable client for accessing remote JSON APIs.
type RPCClient struct {
	// Base is the base URL of the JSON API endpoint.
	Base string

	http http.Client
}

// Call will call an endpoint with with the specified request and store the
// result in the specified response.
func (c *RPCClient) Call(endpoint string, response, request interface{}) error {
	// prepare url
	url := fmt.Sprintf("%s/%s", strings.TrimRight(c.Base, "/"), strings.TrimLeft(endpoint, "/"))

	// handle nil request
	if request == nil {
		request = RPCData(nil)
	}

	// encode request
	buf, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// obtain certificate
	res, err := c.http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}

	// ensure res is closed
	defer res.Body.Close()

	// check code
	if res.StatusCode != 200 {
		var rpcError RPCError
		err = json.NewDecoder(res.Body).Decode(&rpcError)
		if err != nil {
			return err
		}
		return &rpcError
	}

	// decode response if given
	if response != nil {
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			return err
		}
	}

	return nil
}
