package storage_handler

import (
	"net/http"
	"time"

	"objexrg-backend/internal/dtos/bucket/output/response"
	"objexrg-backend/internal/services/storage"
	common_result "objexrg-backend/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type ShareHandler struct {
	service *storage.ShareService
}

func NewShareHandler(service *storage.ShareService) *ShareHandler {
	return &ShareHandler{service: service}
}

type ShareRequest struct {
	ExpiresAt string `json:"expires_at" example:"2026-07-07T23:59:59Z"`
}

// GenerateShareLink godoc
// @Summary Generate a share link for a file
// @Description Generate a public, temporary URL for a file
// @Tags share
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param filename path string true "File name"
// @Param request body ShareRequest false "Optional expiration date (RFC3339 format)"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /share/{bucket}/{filename} [post]
func (h *ShareHandler) Generate(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")
	filename := c.Param("filename")

	if bucketName == "" || filename == "" {
		res := common_result.Fail[any]("bucket and filename are required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	var req ShareRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		res := common_result.Fail[any]("invalid request body", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	var expiresAtPtr *time.Time
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			res := common_result.Fail[any]("invalid expires_at format, must be RFC3339", http.StatusBadRequest)
			c.JSON(res.Code, res)
			return
		}
		expiresAtPtr = &t
	}

	token, err := h.service.GenerateShareLink(userID.(string), bucketName, filename, expiresAtPtr)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	data := map[string]string{
		"token": token,
		"url":   "/s/" + token,
	}

	res := common_result.Ok(&data, http.StatusCreated, "Share link generated successfully")
	c.JSON(res.Code, res)
}

// AccessSharedFile godoc
// @Summary Access a shared file
// @Description Fetch a publicly shared file by its token
// @Tags share
// @Produce octet-stream
// @Param token path string true "Share token"
// @Success 200 {file} binary
// @Failure 404 {object} map[string]string
// @Router /s/{token} [get]
func AccessSharedFile(svc *storage.ShareService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		if token == "" {
			res := common_result.Fail[any]("token required", http.StatusBadRequest)
			c.JSON(res.Code, res)
			return
		}

		stream, obj, err := svc.GetSharedObject(token)
		if err != nil {
			res := common_result.Fail[any]("not found or expired", http.StatusNotFound)
			c.JSON(res.Code, res)
			return
		}
		defer stream.Close()

		serveObject(c, stream, obj.FileName, obj.MimeType, obj.CreatedAt, "inline")
	}
}

// ListShares godoc
// @Summary List all shares
// @Description View all generated shares by the user
// @Tags share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} bucket_response.ShareListResponse
// @Failure 401 {object} map[string]string
// @Router /share/list [get]
func (h *ShareHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	shares, err := h.service.ListUserShares(userID.(string))
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	details := make([]bucket_response.ShareItemResponse, 0)
	for _, s := range shares {
		details = append(details, bucket_response.ShareItemResponse{
			TokenHash:     s.TokenHash,
			Url:           "/s/" + s.TokenHash,
			BucketName:    s.Object.Bucket.Name,
			FileName:      s.Object.FileName,
			MimeType:      s.Object.MimeType,
			SizeBytes:     s.Object.SizeBytes,
			ExpiresAt:     s.ExpiresAt,
			Revoked:       s.Revoked,
			DownloadCount: s.DownloadCount,
			CreatedAt:     s.CreatedAt,
		})
	}

	dataList := bucket_response.ShareListResponse{Details: details}
	res := common_result.Ok(&dataList, http.StatusOK, "Shares retrieved")
	c.JSON(res.Code, res)
}


func (h *ShareHandler) Revoke(c *gin.Context) {
	
}

func (h *ShareHandler) Delete(c *gin.Context) {
	
}

func RegisterShareRoutes(r *gin.RouterGroup, public *gin.RouterGroup, svc *ShareHandler) {
	// Protected route to generate tokens
	r.POST("/share/:bucket/:filename", svc.Generate)
	// Protected route to list tokens
	r.GET("/share/list", svc.List)
	// Public route to access tokens
	public.GET("/s/:token", AccessSharedFile(svc.service))
}
