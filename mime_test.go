package serve

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMimeTypes = map[string]string{
	".html":  "text/html; charset=utf-8",
	".css":   "text/css; charset=utf-8",
	".js":    "application/javascript",
	".wasm":  "application/wasm",
	".jpeg":  "image/jpeg",
	".svg":   "image/svg+xml",
	".ico":   "image/x-icon",
	".mp3":   "audio/mp3",
	".wav":   "audio/x-wav",
	".weba":  "audio/webm",
	".aif":   "audio/x-aiff",
	".flac":  "audio/x-flac",
	".webm":  "video/webm",
	".avi":   "video/x-msvideo",
	".json":  "application/json",
	".csv":   "text/csv; charset=utf-8",
	".gz":    "application/gzip",
	".zip":   "application/zip",
	".tar":   "application/x-tar",
	".pdf":   "application/pdf",
	".woff2": "font/woff2",
}

func init() {
	initMime()
}

func TestMimeTypeByExtension(t *testing.T) {
	for ext, typ := range testMimeTypes {
		assert.Equal(t, typ, MimeTypeByExtension(ext, true), ext)
		assert.Equal(t, typ, MimeTypeByExtension(strings.ToUpper(ext), true), ext)
	}

	assert.Equal(t, 0.0, testing.AllocsPerRun(100, func() {
		MimeTypeByExtension(".html", true)
	}))
}

func TestExtensionsByMimeType(t *testing.T) {
	for ext, typ := range testMimeTypes {
		list, err := ExtensionsByMimeType(typ)
		assert.NoError(t, err)
		assert.Equal(t, ext, list[0])

		list, err = ExtensionsByMimeType(strings.ToUpper(typ))
		assert.NoError(t, err)
		assert.Equal(t, ext, list[0])
	}

	assert.Equal(t, 0.0, testing.AllocsPerRun(100, func() {
		_, err := ExtensionsByMimeType("image/jpeg")
		assert.NoError(t, err)
	}))
}

func BenchmarkMimeTypeByExtension(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MimeTypeByExtension(".html", true)
	}
}

func BenchmarkExtensionsByMimeType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ExtensionsByMimeType("image/jpeg")
		if err != nil {
			panic(err)
		}
	}
}
