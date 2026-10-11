package storage_handler

import (
	"io"
	"net/http"

	bucket_response "objexrg-backend/internal/dtos/bucket/output/response"
	"objexrg-backend/internal/services/storage"
	common_result "objexrg-backend/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type ObjectHandler struct {
	service *storage.ObjectService
}

func NewObjectHandler(service *storage.ObjectService) *ObjectHandler {
	return &ObjectHandler{service: service}
}

type setVisibilityRequest struct {
	IsPublic *bool `json:"is_public"`
}

// SetVisibility godoc
// @Summary Change object visibility
// @Description Toggle whether an object is publicly readable
// @Tags objects
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param filename path string true "File name"
// @Param request body object.VisibilityRequest true "Visibility payload"
// @Security BearerAuth
// @Success 200 {object} bucket_response.ObjectVisibilityResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /bucket/{bucket}/objects/{filename} [patch]
func (h *ObjectHandler) SetVisibility(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")
	filename := c.Param("filename")

	var req setVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.IsPublic == nil {
		res := common_result.Fail[any]("is_public (boolean) is required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	obj, err := h.service.SetVisibility(userID.(string), bucketName, filename, *req.IsPublic)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}

	data := bucket_response.ObjectVisibilityResponse{
		ID:       obj.ID.String(),
		Bucket:   bucketName,
		FileName: obj.FileName,
		IsPublic: obj.IsPublic,
	}
	if obj.IsPublic {
		data.PublicURL = "/public/" + obj.ID.String()
	}

	res := common_result.Ok(&data, http.StatusOK, "Visibility updated successfully")
	c.JSON(res.Code, res)
}

// DeleteObject godoc
// @Summary Delete an object
// @Description Delete a file from a bucket
// @Tags objects
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param filename path string true "File name"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /bucket/{bucket}/objects/{filename} [delete]
func (h *ObjectHandler) DeleteObject(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")
	filename := c.Param("filename")

	if err := h.service.DeleteObject(userID.(string), bucketName, filename); err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}

	res := common_result.Ok[any](nil, http.StatusOK, "File deleted successfully")
	c.JSON(res.Code, res)
}

// GetPublicObject godoc
// @Summary Read a public object
// @Description Stream a publicly shared object (no authentication required)
// @Tags public
// @Produce octet-stream
// @Param id path string true "Object ID"
// @Success 200 {file} binary
// @Failure 404 {object} map[string]string
// @Router /public/{id} [get]
func (h *ObjectHandler) GetPublicObject(c *gin.Context) {
	id := c.Param("id")

	stream, obj, err := h.service.GetPublic(id)
	if err != nil {
		res := common_result.Fail[any]("not found", http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}
	defer stream.Close()

	c.Header("Content-Type", obj.MimeType)
	c.Header("Content-Disposition", "inline; filename=\""+obj.FileName+"\"")
	c.Status(http.StatusOK)

	_, _ = io.Copy(c.Writer, stream)
}

func RegisterObjectRoutes(r *gin.RouterGroup, h *ObjectHandler) {
	r.PATCH("/bucket/:bucket/objects/:filename", h.SetVisibility)
	r.DELETE("/bucket/:bucket/objects/:filename", h.DeleteObject)
}

func RegisterPublicObjectRoutes(public *gin.RouterGroup, h *ObjectHandler) {
	public.GET("/public/:id", h.GetPublicObject)
}
