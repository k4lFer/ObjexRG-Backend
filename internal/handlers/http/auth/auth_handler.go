package auth_handler

import (
	"net/http"

	"objexrg/internal/dtos/auth/input"
	auth_response "objexrg/internal/dtos/auth/response"
	authsvc "objexrg/internal/services/auth"
	common_result "objexrg/internal/shared/common"
	authvalidators "objexrg/internal/validators/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService       *authsvc.AuthService
	registerValidator *authvalidators.RegisterValidator
	loginValidator    *authvalidators.LoginValidator
}

func NewAuthHandler(authService *authsvc.AuthService, registerValidator *authvalidators.RegisterValidator, loginValidator *authvalidators.LoginValidator) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		registerValidator: registerValidator,
		loginValidator:    loginValidator,
	}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body auth_input.RegisterRequest true "User registration details"
// @Success 201 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth_input.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		res := common_result.Fail[any]("Invalid JSON format", http.StatusUnprocessableEntity)
		c.JSON(res.Code, res)
		return
	}

	result := h.registerValidator.Validate(&req)
	if !result.IsValid {
		res := common_result.Fail[any](result.Message, result.StatusCode)
		c.JSON(res.Code, res)
		return
	}

	if err := h.authService.Register(&req); err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	res := common_result.Ok[any](nil, http.StatusCreated, "User registered successfully")
	c.JSON(res.Code, res)
}

/*
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth_input.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Validar con validator
	result := h.registerValidator.Validate(&req)
	if !result.IsValid {
		c.JSON(result.StatusCode, gin.H{"error": result.Message})
		return
	}

	if err := h.authService.Register(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

*/
// LoginUser godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth_input.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth_input.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		res := common_result.Fail[any]("Invalid JSON format", http.StatusUnprocessableEntity)
		c.JSON(res.Code, res)
		return
	}

	result := h.loginValidator.Validate(&req)
	if !result.IsValid {
		res := common_result.Fail[any](result.Message, result.StatusCode)
		c.JSON(res.Code, res)
		return
	}

	user, tokens, err := h.authService.Login(&req)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}

	data := auth_response.LoginResponse{
		ID:           user.ID.String(),
		Username:     user.Username,
		FirstName:    user.UserData.FirstName,
		LastName:     user.UserData.LastName,
		Email:        user.Email,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	}

	res := common_result.Ok(&data, http.StatusOK, "Login successful")
	c.JSON(res.Code, res)
}



/*
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth_input.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Validar con validator
	result := h.loginValidator.Validate(&req)
	if !result.IsValid {
		c.JSON(result.StatusCode, gin.H{"error": result.Message})
		return
	}

	user, tokens, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Devolver datos del usuario y tokens
	response := gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"first_name":    user.UserData.FirstName,
		"last_name":     user.UserData.LastName,
		"email":         user.Email,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}
*/

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth_input.RefreshRequest true "Refresh token details"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth_input.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := common_result.Fail[any]("Invalid JSON format", http.StatusUnprocessableEntity)
		c.JSON(res.Code, res)
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusUnauthorized)
		c.JSON(res.Code, res)
		return
	}
	data := auth_response.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
	}

	res := common_result.Ok(&data, http.StatusOK, "Token refreshed successfully")
	c.JSON(res.Code, res)
}

// LogoutUser godoc
// @Summary Logout user
// @Description Revoke the provided refresh token so it can no longer be used
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth_input.LogoutRequest true "Refresh token to revoke"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req auth_input.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := common_result.Fail[any]("Invalid JSON format", http.StatusUnprocessableEntity)
		c.JSON(res.Code, res)
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		res := common_result.Fail[any](err.Error(), http.StatusInternalServerError)
		c.JSON(res.Code, res)
		return
	}

	res := common_result.Ok[any](nil, http.StatusOK, "Logged out successfully")
	c.JSON(res.Code, res)
}

func RegisterAuthRoutes(router *gin.RouterGroup, authHandler *AuthHandler) {
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/refresh", authHandler.RefreshToken)
	router.POST("/logout", authHandler.Logout)
}
