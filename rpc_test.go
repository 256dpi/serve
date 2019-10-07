package serve

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRPC(t *testing.T) {
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    "0.0.0.0:1337",
		Handler: mux,
	}

	go server.ListenAndServe()
	defer server.Close()
	time.Sleep(10 * time.Millisecond)

	client := RPCClient{
		Base: "http://0.0.0.0:1337",
	}

	/* nil response */

	mux.Handle("/nil", RPCHandler(0, func(ctx *RPCContext) interface{} {
		return nil
	}))

	var res RPCData
	err := client.Call("/nil", &res, nil)
	assert.NoError(t, err)
	assert.Equal(t, RPCData(nil), res)

	/* arbitrary data */

	mux.Handle("/data", RPCHandler(0, func(ctx *RPCContext) interface{} {
		return RPCData{"foo": 42}
	}))

	err = client.Call("/data", &res, nil)
	assert.NoError(t, err)
	assert.Equal(t, RPCData{"foo": float64(42)}, res)

	/* basic types */

	mux.Handle("/int", RPCHandler(0, func(ctx *RPCContext) interface{} {
		return 42
	}))

	var i int
	err = client.Call("/int", &i, nil)
	assert.NoError(t, err)
	assert.Equal(t, 42, i)

	/* structures */

	type Item struct {
		Name string `json:"foo"`
	}

	mux.Handle("/item", RPCHandler(100, func(ctx *RPCContext) interface{} {
		var item Item
		return ctx.Handle(&item, func() interface{} {
			return Item{Name: strings.ToUpper(item.Name)}
		})
	}))

	var item Item
	err = client.Call("/item", &item, Item{Name: "foo"})
	assert.NoError(t, err)
	assert.Equal(t, Item{Name: "FOO"}, item)

	/* raw errors */

	mux.Handle("/fail", RPCHandler(0, func(ctx *RPCContext) interface{} {
		return fmt.Errorf("some error")
	}))

	err = client.Call("/fail", nil, nil)
	assert.Equal(t, &RPCError{
		Status: 500,
		Title:  "internal server error",
		Detail: "unknown error",
	}, err)

	/* extended errors */

	mux.Handle("/err", RPCHandler(0, func(ctx *RPCContext) interface{} {
		return RPCBadRequest("just bad", "param.foo")
	}))

	err = client.Call("/err", nil, nil)
	assert.Equal(t, &RPCError{
		Status: 400,
		Title:  "bad request",
		Detail: "just bad",
		Source: "param.foo",
	}, err)
}
