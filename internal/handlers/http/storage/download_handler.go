package storage_handler

import (
	"net/http"

	"objexrg/internal/services/storage"
	common_result "objexrg/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	service  *storage.ObjectService
	delivery *Delivery
}

func NewDownloadHandler(service *storage.ObjectService, delivery *Delivery) *DownloadHandler {
	return &DownloadHandler{service: service, delivery: delivery}
}

// DownloadFile godoc
// @Summary Download a file from a bucket
// @Description Download the specified file from the bucket
// @Tags download
// @Accept json
// @Produce octet-stream
// @Param bucket path string true "Bucket name"
// @Param filename path string true "File name"
// @Security BearerAuth
// @Success 200 {file} binary
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /download/{bucket}/{filename} [get]
func (h *DownloadHandler) Download(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")
	filename := c.Param("filename")
	if bucketName == "" || filename == "" {
		res := common_result.Fail[any]("bucket and filename required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	obj, err := h.service.DownloadMeta(c.Request.Context(), userID.(string), bucketName, filename)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}

	h.delivery.Deliver(c, obj.OwnerID.String(), bucketName, obj, "attachment")
}

func RegisterDownloadRoutes(r *gin.RouterGroup, h *DownloadHandler) {
	r.GET("/:bucket/:filename", h.Download)
}
