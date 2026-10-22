package storage_handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"objexrg/internal/dtos/bucket/output/response"
	"objexrg/internal/services/storage"
	common_result "objexrg/internal/shared/common"

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
// @Description Upload a file to the specified bucket (proxy path)
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

	key, err := h.service.Upload(c.Request.Context(), userID.(string), bucketName, fileHeader.Filename, fileHeader.Size, mimeType, body)
	if err != nil {
		h.writeUploadError(c, err)
		return
	}

	data := response.UploadDataResponse{
		Bucket:   bucketName,
		Filename: fileHeader.Filename,
		Path:     key,
		Size:     fileHeader.Size,
	}

	res := common_result.Ok(&data, http.StatusCreated, "File uploaded successfully")
	c.JSON(res.Code, res)
}

type presignUploadRequest struct {
	FileName  string `json:"fileName"`
	SizeBytes int64  `json:"sizeBytes"`
	MimeType  string `json:"mimeType"`
}

// PresignUpload godoc
// @Summary Request a presigned upload URL
// @Description Reserve a storage key and return a direct upload URL when the driver supports it
// @Tags upload
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param request body presignUploadRequest true "Upload intent"
// @Security BearerAuth
// @Success 200 {object} response.PresignUploadResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Router /upload/{bucket}/presign [post]
func (h *UploadHandler) PresignUpload(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")

	var req presignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.FileName == "" {
		res := common_result.Fail[any]("fileName is required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	result, err := h.service.PresignUpload(c.Request.Context(), userID.(string), bucketName, req.FileName, req.SizeBytes, req.MimeType)
	if err != nil {
		h.writeUploadError(c, err)
		return
	}

	data := response.PresignUploadResponse{
		Mode:       result.Mode,
		StorageKey: result.StorageKey,
		Key:        result.Key,
	}
	switch result.Mode {
	case "s3-post":
		data.Method = "POST"
		data.UploadURL = result.UploadURL
		data.Fields = result.Fields
	case "s3":
		data.Method = "PUT"
		data.UploadURL = result.UploadURL
		data.Headers = result.Headers
	}

	res := common_result.Ok(&data, http.StatusOK, "Upload URL generated successfully")
	c.JSON(res.Code, res)
}

type finalizeUploadRequest struct {
	FileName   string `json:"fileName"`
	StorageKey string `json:"storageKey"`
}

// FinalizeUpload godoc
// @Summary Finalize a direct upload
// @Description Validate and register an object uploaded directly to storage
// @Tags upload
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param request body finalizeUploadRequest true "Finalize payload"
// @Security BearerAuth
// @Success 201 {object} response.ObjectUploadResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Router /upload/{bucket}/finalize [post]
func (h *UploadHandler) FinalizeUpload(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")

	var req finalizeUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.FileName == "" || req.StorageKey == "" {
		res := common_result.Fail[any]("fileName and storageKey are required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	obj, err := h.service.FinalizeUpload(c.Request.Context(), userID.(string), bucketName, req.FileName, req.StorageKey)
	if err != nil {
		h.writeUploadError(c, err)
		return
	}

	data := response.ObjectUploadResponse{
		ID:         obj.ID.String(),
		Bucket:     bucketName,
		FileName:   obj.FileName,
		SizeBytes:  obj.SizeBytes,
		MimeType:   obj.MimeType,
		StorageKey: obj.StorageKey,
		Key:        obj.Path,
		CreatedAt:  obj.CreatedAt,
	}

	res := common_result.Ok(&data, http.StatusCreated, "File uploaded successfully")
	c.JSON(res.Code, res)
}

func (h *UploadHandler) writeUploadError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, storage.ErrObjectExists):
		res := common_result.Fail[any](err.Error(), http.StatusConflict)
		c.JSON(res.Code, res)
	case errors.Is(err, storage.ErrUploadTooLarge):
		res := common_result.Fail[any](h.tooLargeMessage(), http.StatusRequestEntityTooLarge)
		c.JSON(res.Code, res)
	default:
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
	}
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
	r.POST("/:bucket/presign", svc.PresignUpload)
	r.POST("/:bucket/finalize", svc.FinalizeUpload)
}
