package storage_handler

import (
	"io"
	"net/http"

	"objexrg-backend/internal/services/storage"
	common_result "objexrg-backend/internal/shared/common"

	"github.com/gin-gonic/gin"
)

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
func Download(svc *storage.UploadService) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		bucketPath := userID.(string) + "/" + bucketName
		stream, err := svc.Download(bucketPath, filename)
		if err != nil {
			res := common_result.Fail[any](err.Error(), http.StatusNotFound)
			c.JSON(res.Code, res)
			return
		}
		defer stream.Close()

		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Type", "application/octet-stream")
		c.Status(http.StatusOK)

		_, _ = io.Copy(c.Writer, stream)
	}
}

func RegisterDownloadRoutes(r *gin.RouterGroup, svc *storage.UploadService) {
	r.GET("/:bucket/:filename", Download(svc))
}
