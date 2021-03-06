package serve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostname(t *testing.T) {
	matrix := []struct {
		i string
		o string
	}{
		{
			i: "foo.com",
			o: "foo.com",
		},
		{
			i: "foo.com:80",
			o: "foo.com",
		},
		{
			i: "foo.com:bar",
			o: "foo.com:bar",
		},
	}

	for _, item := range matrix {
		assert.Equal(t, item.o, Hostname(item.i))
	}
}

func TestIP(t *testing.T) {
	matrix := []struct {
		i string
		o string
	}{
		{
			i: "1.1.1.1:80",
			o: "1.1.1.1",
		},
		{
			i: "1.1.1.1",
			o: "1.1.1.1",
		},
		{
			i: "[::1]:80",
			o: "::1",
		},
		{
			i: "::1",
			o: "::1",
		},
		{
			i: ":80",
			o: "",
		},
		{
			i: "foo",
			o: "foo",
		},
	}

	for _, item := range matrix {
		assert.Equal(t, item.o, IP(item.i), item)
	}
}
