package user_handler

import (
	"net/http"
	"objexrg-backend/internal/dtos/user/output/response"
	"objexrg-backend/internal/dtos/user/input"
	user_service "objexrg-backend/internal/services/user"
	common_result "objexrg-backend/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *user_service.UserService
}

func NewUserHandler(userService *user_service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}


func (h *UserHandler) GetById(c *gin.Context) {
	
}

// Me godoc
// @Summary Get current user info
// @Description Get the current authenticated user's information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user_response.MeResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		r := common_result.Fail[any]("User ID not found in context", http.StatusUnauthorized)
		c.JSON(r.Code, r)
		return
	}
	
	user, err := h.userService.GetById(userID.(string))
	if err != nil {
		r := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(r.Code, r)
		return
	}

	data := user_response.MeResponse {
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.UserData.FirstName,
		LastName:  user.UserData.LastName,
		BirthDate: user.UserData.BirthDate,
		CreatedAt: user.UserData.CreatedAt,
	}
	res := common_result.Ok(&data, http.StatusOK, "Your profile was retrieved successfully")
	c.JSON(res.Code, res)
}

// PatchMe godoc
// @Summary Update current user info
// @Description Update the current authenticated user's information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body user_input.PatchMeRequest true "Updated user info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/me [patch]
func (h *UserHandler) PatchMe(c *gin.Context) {
	var req user_input.PatchMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := common_result.Fail[any]("Invalid JSON format", http.StatusBadRequest)
		c.JSON(res.Code, res)
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		res := common_result.Fail[any]("User ID not found in context", http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}
	res := common_result.Ok(&userID, http.StatusOK, "User updated successfully")
	c.JSON(res.Code, res)
}

func RegisterUserRoutes(router *gin.RouterGroup, userHandler *UserHandler) {
	router.GET("/me", userHandler.Me)
	router.PATCH("/me", userHandler.PatchMe)
}