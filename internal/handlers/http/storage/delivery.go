package storage_handler

import (
	"errors"
	"net/http"
	"time"

	"objexrg/internal/interfaces"
	storage_models "objexrg/internal/models/storage"
	common_result "objexrg/internal/shared/common"

	"github.com/gin-gonic/gin"
)

// Delivery resolves how a stored object is handed to the client: a 302
// redirect to a presigned URL on S3 drivers, or a proxied stream with Range
// support on the local driver.
type Delivery struct {
	store interfaces.BlobStore
	ttl   time.Duration
}

func NewDelivery(store interfaces.BlobStore, ttl time.Duration) *Delivery {
	return &Delivery{store: store, ttl: ttl}
}

func (d *Delivery) Deliver(c *gin.Context, ownerID, bucketName string, obj *storage_models.Objects, disposition string) {
	key := interfaces.ObjectKey(ownerID, bucketName, obj.StorageName())

	url, err := d.store.PresignGetURL(c.Request.Context(), key, obj.MimeType, disposition, obj.FileName, d.ttl)
	if err == nil {
		c.Redirect(http.StatusFound, url)
		return
	}
	if !errors.Is(err, interfaces.ErrPresignNotSupported) {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	stream, err := d.store.Get(c.Request.Context(), key)
	if err != nil {
		res := common_result.Fail[any]("not found", http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}
	defer stream.Close()

	serveObject(c, stream, obj.FileName, obj.MimeType, obj.CreatedAt, disposition)
}
