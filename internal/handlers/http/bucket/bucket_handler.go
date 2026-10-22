package bucket_handler

import (
	"net/http"

	"objexrg/internal/dtos/bucket/input"
	"objexrg/internal/dtos/bucket/output/response"
	bucketsvc "objexrg/internal/services/bucket"
	common_result "objexrg/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type BucketHandler struct {
	svc *bucketsvc.BucketService
}

func NewBucketHandler(svc *bucketsvc.BucketService) *BucketHandler {
	return &BucketHandler{svc: svc}
}

// CreateBucket godoc
// @Summary Create a new bucket
// @Description Create a new bucket with the given name
// @Tags buckets
// @Accept json
// @Produce json
// @Param bucket body bucket_input.RegisterRequest true "Bucket registration details"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bucket [post]
func (h *BucketHandler) CreateBucket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	var req bucket_input.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	if err := h.svc.CreateSubBucket(userID.(string), req.Name); err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	data := response.CreateBucketResponse{
		Bucket: req.Name,
		User:   userID.(string),
	}

	res := common_result.Ok(&data, http.StatusCreated, "Bucket created successfully")
	c.JSON(res.Code, res)
}


/*
func (h *BucketHandler) CreateBucket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req bucket_input.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.CreateSubBucket(userID.(string), req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"bucket": req.Name, "user": userID, "status": "created"})
}*/

// DeleteBucket godoc
// @Summary Delete a bucket
// @Description Delete the bucket with the given name
// @Tags buckets
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bucket/{bucket} [delete]
func (h *BucketHandler) DeleteBucket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	name := c.Param("bucket")
	if name == "" {
		res := common_result.Fail[any]("bucket name required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}
	if err := h.svc.DeleteSubBucket(userID.(string), name); err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}
	res := common_result.Ok[any](nil, http.StatusOK, "Bucket deleted successfully")
	c.JSON(res.Code, res)
}

// ListBuckets godoc
// @Summary List all buckets
// @Description Get a list of all buckets
// @Tags buckets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /bucket [get]
func (h *BucketHandler) ListBuckets(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	buckets, err := h.svc.ListBuckets(userID.(string))
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	data := response.BucketsListResponse {
		Buckets: buckets,
	}

	res := common_result.Ok(&data, http.StatusOK, "Buckets retrieved successfully")
	c.JSON(res.Code, res)
}

// ListObjects godoc
// @Summary List objects in a bucket
// @Description Get a list of all objects in the specified bucket
// @Tags buckets
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bucket/{bucket} [get]
func (h *BucketHandler) ListObjects(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	name := c.Param("bucket")
	if name == "" {
		res := common_result.Fail[any]("bucket name required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}

	objects, err := h.svc.ListObjects(userID.(string), name)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	data := response.ObjectsListResponse{
		Bucket:  name,
		Objects: objects,
	}

	res := common_result.Ok(&data, http.StatusOK, "Objects retrieved successfully")
	c.JSON(res.Code, res)
}

/*func (h *BucketHandler) ListObjects(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	name := c.Param("bucket")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name required"})
		return
	}
	bucketPath := userID.(string) + "/" + name
	objects, err := h.storage.ListObjects(bucketPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bucket": name, "objects": objects})
}
*/
// ListStorage godoc
// @Summary List storage in a bucket
// @Description Get a list of all storage objects in the specified bucket
// @Tags buckets
// @Accept json
// @Produce json
// @Param bucket path string true "Bucket name"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bucket/{bucket}/storage [get]
func (h *BucketHandler) ListStorage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("unauthorized", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	name := c.Param("bucket")
	if name == "" {
		res := common_result.Fail[any]("bucket name required", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}
	objects, err := h.svc.ListObjects(userID.(string), name)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	data := response.ObjectsListResponse{
		Bucket:  name,
		Objects: objects,
	}

	res := common_result.Ok(&data, http.StatusOK, "Storage retrieved successfully")
	c.JSON(res.Code, res)
}

func RegisterBucketRoutes(r *gin.RouterGroup, svc *bucketsvc.BucketService) {
	handler := NewBucketHandler(svc)
	r.POST("", handler.CreateBucket)
	r.DELETE("/:bucket", handler.DeleteBucket)
	r.GET("", handler.ListBuckets)
	r.GET("/:bucket", handler.ListObjects)
	r.GET("/:bucket/storage", handler.ListStorage)
}
