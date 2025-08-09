package pathx

import (
	"path/filepath"
	"regexp"
	"strings"
)

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9-_]`)

// SafeName normalizes a string to a safe directory name.
func SafeName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return nonAlnum.ReplaceAllString(s, "-")
}

// Join is a thin wrapper to filepath.Join.
func Join(elem ...string) string { return filepath.Join(elem...) }
