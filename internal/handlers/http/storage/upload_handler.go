package storage_handler

import (
	"net/http"

	bucket_response "objexrg-backend/internal/dtos/bucket/output/response"
	"objexrg-backend/internal/services/storage"
	common_result "objexrg-backend/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *storage.UploadService
}

func NewUploadHandler(service *storage.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

// UploadFile godoc
// @Summary Upload a file to a bucket
// @Description Upload a file to the specified bucket
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param file formData file true "File to upload"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /upload/{bucket} [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")
	if bucketName == "" {
		res := common_result.Fail[any]("bucket required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	// Verificar que el bucket existe
	if !h.service.BucketExists(userID.(string), bucketName) {
		res := common_result.Fail[any]("bucket does not exist", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		res := common_result.Fail[any]("file is required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}
	defer src.Close()

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	stored, err := h.service.Upload(userID.(string), bucketName, fileHeader.Filename, fileHeader.Size, mimeType, src)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}
	
	data := bucket_response.UploadDataResponse{
		Bucket:   bucketName,
		Filename: fileHeader.Filename,
		Path:     stored,
		Size:     fileHeader.Size,
	}

	res := common_result.Ok(&data, http.StatusCreated, "File uploaded successfully")
	c.JSON(res.Code, res)
}

func RegisterUploadRoutes(r *gin.RouterGroup, svc *UploadHandler) {
	r.POST("/:bucket", svc.Upload)
}
