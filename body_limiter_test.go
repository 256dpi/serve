package serve

import (
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ io.ReadCloser = &BodyLimiter{}

func TestLimitBody(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.org", strings.NewReader("hello world"))
	w := httptest.NewRecorder()

	orig := r.Body

	LimitBody(w, r, 2)
	assert.Equal(t, orig, r.Body.(*BodyLimiter).Original)

	LimitBody(w, r, 5)
	assert.Equal(t, orig, r.Body.(*BodyLimiter).Original)

	bytes, err := ioutil.ReadAll(r.Body)
	assert.Error(t, err)
	assert.Equal(t, "hello", string(bytes))
	assert.Equal(t, err, ErrBodyLimitExceeded)
}

func TestDataSize(t *testing.T) {
	assert.Equal(t, uint64(50*1000), DataSize("50K"))
	assert.Equal(t, uint64(5*1000*1000), DataSize("5M"))
	assert.Equal(t, uint64(100*1000*1000*1000), DataSize("100G"))

	for _, str := range []string{"", "1", "K", "10", "KM"} {
		assert.PanicsWithValue(t, `fire: data size must be like 4K, 20M or 5G`, func() {
			DataSize(str)
		})
	}
}
