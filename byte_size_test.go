package serve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteSize(t *testing.T) {
	assert.Equal(t, int64(50*1000), MustByteSize("50K"))
	assert.Equal(t, int64(5*1000*1000), MustByteSize("5MiB"))
	assert.Equal(t, int64(100*1024*1024*1024), MustByteSize("100GB"))

	for _, str := range []string{"", "1", "K", "10", "KM"} {
		assert.PanicsWithValue(t, ErrInvalidByteSize, func() {
			MustByteSize(str)
		})
	}
}
