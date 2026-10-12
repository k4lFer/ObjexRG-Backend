package storage_handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const fallbackMime = "application/octet-stream"

// serveObject streams an object to the client honoring HTTP Range and HEAD
// requests when the underlying stream supports seeking (local files do).
func serveObject(
	c *gin.Context,
	stream io.ReadCloser,
	name string,
	mimeType string,
	modTime time.Time,
	disposition string,
) {
	if mimeType == "" {
		mimeType = fallbackMime
	}

	c.Header("Content-Type", mimeType)
	c.Header("Content-Disposition", disposition+`; filename="`+name+`"`)

	if seeker, ok := stream.(io.ReadSeeker); ok {
		// ServeContent handles Range, HEAD, Content-Length and 206 responses.
		http.ServeContent(c.Writer, c.Request, name, modTime, seeker)
		return
	}

	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, stream)
}
