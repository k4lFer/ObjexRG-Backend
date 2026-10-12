package storage_handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	bucket_response "objexrg-backend/internal/dtos/bucket/output/response"
	"objexrg-backend/internal/services/storage"
	common_result "objexrg-backend/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service        *storage.UploadService
	maxUploadBytes int64
}

func NewUploadHandler(service *storage.UploadService, maxUploadBytes int64) *UploadHandler {
	return &UploadHandler{service: service, maxUploadBytes: maxUploadBytes}
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
// @Failure 409 {object} map[string]string
// @Failure 413 {object} map[string]string
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

	// Hard cap on the request body so oversized uploads are rejected early.
	if h.maxUploadBytes > 0 {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxUploadBytes)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			res := common_result.Fail[any](h.tooLargeMessage(), http.StatusRequestEntityTooLarge)
			c.JSON(res.Code, res)
			return
		}
		res := common_result.Fail[any]("file is required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	if h.maxUploadBytes > 0 && fileHeader.Size > h.maxUploadBytes {
		res := common_result.Fail[any](h.tooLargeMessage(), http.StatusRequestEntityTooLarge)
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

	// Sniff the content type from the leading bytes, then prepend them back so
	// the full stream is preserved for storage.
	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	head = head[:n]
	mimeType := chooseMimeType(head, fileHeader.Header.Get("Content-Type"))
	body := io.MultiReader(bytes.NewReader(head), src)

	stored, err := h.service.Upload(userID.(string), bucketName, fileHeader.Filename, fileHeader.Size, mimeType, body)
	if err != nil {
		if errors.Is(err, storage.ErrObjectExists) {
			res := common_result.Fail[any](err.Error(), http.StatusConflict)
			c.JSON(res.Code, res)
			return
		}
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

func (h *UploadHandler) tooLargeMessage() string {
	if h.maxUploadBytes <= 0 {
		return "file is too large"
	}
	mib := float64(h.maxUploadBytes) / (1024 * 1024)
	return "file exceeds the maximum allowed size of " +
		strconv.FormatFloat(mib, 'f', -1, 64) + " MiB"
}

// chooseMimeType prefers the sniffed type, falling back to the declared one
// when sniffing yields the generic octet-stream (common for videos).
func chooseMimeType(head []byte, declared string) string {
	declared = strings.TrimSpace(declared)

	sniffed := http.DetectContentType(head)
	if sniffed != "application/octet-stream" {
		return sniffed
	}
	if declared != "" {
		return declared
	}
	return "application/octet-stream"
}

func RegisterUploadRoutes(r *gin.RouterGroup, svc *UploadHandler) {
	r.POST("/:bucket", svc.Upload)
}
