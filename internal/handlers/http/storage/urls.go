package storage_handler

import "strings"

// absoluteURL joins the configured public base URL with a path so clients get
// ready-to-share links regardless of where they run. When no base is set the
// path is returned as-is (relative).
func absoluteURL(base, path string) string {
	base = strings.TrimRight(base, "/")
	if base == "" {
		return path
	}
	return base + path
}
