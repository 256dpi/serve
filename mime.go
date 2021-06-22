package serve

import (
	// for embeds
	_ "embed"
	"encoding/json"
	"mime"
	"sort"
	"strings"
	"sync"
)

//go:embed mime.json
var mimeJSON []byte

var mimeOnce sync.Once
var mimeDB = map[string]mimeEntry{}
var mimeExt = map[string][]mimeEntry{}

var mimeSources = map[string]int{
	"nginx":  0,
	"apache": 1,
	"":       2,
	"iana":   3,
}

type mimeEntry struct {
	Name            string   `json:"-"`
	NameWithCharset string   `json:"-"`
	Source          string   `json:"source"`
	Compressible    bool     `json:"compressible"`
	Extensions      []string `json:"extensions"`
}

func initMime() {
	mimeOnce.Do(func() {
		// decode database
		err := json.Unmarshal(mimeJSON, &mimeDB)
		if err != nil {
			panic(err)
		}

		// post-process entries
		for name, entry := range mimeDB {
			// set name
			entry.Name = name

			// add name with charset
			if strings.HasPrefix(name, "text/") {
				entry.NameWithCharset = name + "; charset=utf-8"
			}

			// add dot prefix
			for i, ext := range entry.Extensions {
				entry.Extensions[i] = "." + ext
			}

			// store by extensions
			for _, ext := range entry.Extensions {
				mimeExt[ext] = append(mimeExt[ext], entry)
			}
		}

		// sort ext entries
		for _, entries := range mimeExt {
			sort.Slice(entries, func(i, j int) bool {
				return mimeSources[entries[i].Source] < mimeSources[entries[j].Source]
			})
		}
	})
}

// MimeTypeByExtension returns the MIME type associated with the provided file
// extension. The extension ext should begin with a leading dot. When ext has no
// associated type, it returns "".
//
// Note: It will prefer a static DB over the to the builtin mime package.
func MimeTypeByExtension(ext string, withCharset bool) string {
	// initialize
	initMime()

	// lower case
	ext = strings.ToLower(ext)

	// check table
	var name string
	entries, ok := mimeExt[ext]
	if ok {
		name = entries[0].Name
		if withCharset && entries[0].NameWithCharset != "" {
			name = entries[0].NameWithCharset
		}
	} else {
		name = mime.TypeByExtension(ext)
		if withCharset && strings.HasPrefix(name, "text/") && !strings.Contains(name, "charset=") {
			name += "; charset=utf-8"
		}
	}

	return name
}

// ExtensionsByMimeType returns the extensions known to be associated with the
// provided MIME type. The returned extensions will each begin with a leading dot.
// When typ has no associated extensions, it returns a nil slice.
//
// Note: It will prefer a static DB over the to the builtin mime package.
func ExtensionsByMimeType(typ string) ([]string, error) {
	// initialize
	initMime()

	// get key
	var key string
	if !strings.ContainsRune(typ, ';') {
		key = strings.ToLower(typ)
	} else {
		var err error
		key, _, err = mime.ParseMediaType(typ)
		if err != nil {
			return nil, err
		}
	}

	// check table
	entry, ok := mimeDB[key]
	if !ok {
		return mime.ExtensionsByType(typ)
	}

	return entry.Extensions, nil
}
