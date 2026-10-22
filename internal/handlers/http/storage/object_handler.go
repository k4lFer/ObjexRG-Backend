package storage_handler

import (
	"net/http"

	"objexrg/internal/dtos/bucket/output/response"
	"objexrg/internal/services/storage"
	common_result "objexrg/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type ObjectHandler struct {
	service       *storage.ObjectService
	delivery      *Delivery
	publicBaseURL string
}

func NewObjectHandler(service *storage.ObjectService, delivery *Delivery, publicBaseURL string) *ObjectHandler {
	return &ObjectHandler{service: service, delivery: delivery, publicBaseURL: publicBaseURL}
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
// @Param request body map[string]interface{} true "Visibility payload (is_public)"
// @Security BearerAuth
// @Success 200 {object} response.ObjectVisibilityResponse
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

	data := response.ObjectVisibilityResponse{
		ID:       obj.ID.String(),
		Bucket:   bucketName,
		FileName: obj.FileName,
		IsPublic: obj.IsPublic,
	}
	if obj.IsPublic {
		data.PublicURL = absoluteURL(h.publicBaseURL, "/f/"+obj.ID.String())
	}

	res := common_result.Ok(&data, http.StatusOK, "Visibility updated successfully")
	c.JSON(res.Code, res)
}

// Link godoc
// @Summary Get a direct URL for an object
// @Description Return a presigned URL (S3) or the proxy download path (local) for an object
// @Tags objects
// @Produce json
// @Param bucket path string true "Bucket name"
// @Param filename path string true "File name"
// @Param disposition query string false "inline or attachment"
// @Security BearerAuth
// @Success 200 {object} response.ObjectLinkResponse
// @Failure 404 {object} map[string]string
// @Router /bucket/{bucket}/objects/{filename}/link [get]
func (h *ObjectHandler) Link(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	bucketName := c.Param("bucket")
	filename := c.Param("filename")

	disposition := c.DefaultQuery("disposition", "inline")
	if disposition != "attachment" {
		disposition = "inline"
	}

	url, err := h.service.LinkURL(c.Request.Context(), userID.(string), bucketName, filename, disposition)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}

	data := response.ObjectLinkResponse{
		URL:         url,
		FileName:    filename,
		Disposition: disposition,
	}

	res := common_result.Ok(&data, http.StatusOK, "Link generated successfully")
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

	if err := h.service.DeleteObject(c.Request.Context(), userID.(string), bucketName, filename); err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}

	res := common_result.Ok[any](nil, http.StatusOK, "File deleted successfully")
	c.JSON(res.Code, res)
}

// GetPublicObject godoc
// @Summary Read a public object
// @Description Stream (or redirect to) a publicly shared object (no authentication required)
// @Tags public
// @Produce octet-stream
// @Param id path string true "Object ID"
// @Success 200 {file} binary
// @Success 302 {string} string "Redirect to presigned URL"
// @Failure 404 {object} map[string]string
// @Router /f/{id} [get]
func (h *ObjectHandler) GetPublicObject(c *gin.Context) {
	id := c.Param("id")

	obj, err := h.service.GetPublic(c.Request.Context(), id)
	if err != nil {
		res := common_result.Fail[any]("not found", http.StatusNotFound)
		c.JSON(res.Code, res)
		return
	}

	h.delivery.Deliver(c, obj.OwnerID.String(), obj.Bucket.Name, obj, "inline")
}

func RegisterObjectRoutes(r *gin.RouterGroup, h *ObjectHandler) {
	r.PATCH("/bucket/:bucket/objects/:filename", h.SetVisibility)
	r.DELETE("/bucket/:bucket/objects/:filename", h.DeleteObject)
	r.GET("/bucket/:bucket/objects/:filename/link", h.Link)
}

func RegisterPublicObjectRoutes(public *gin.RouterGroup, h *ObjectHandler) {
	public.GET("/f/:id", h.GetPublicObject)
}
