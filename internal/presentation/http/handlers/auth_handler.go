package handlers

import (
	"net/http"
	"strings"
	"usermanagement-api/internal/application/dto"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/constants"
	"usermanagement-api/internal/presentation/http/middleware"
	"usermanagement-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	// Accepts both application/x-www-form-urlencoded (FastAPI OAuth2PasswordRequestForm)
	// and JSON bodies; username+password fields.
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(err.Error(), "VALIDATION_ERROR", nil))
		return
	}

	resp, err := h.authUseCase.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(err.Error(), "UNAUTHORIZED", nil))
		return
	}

	// Mirror FastAPI: AuthToken header + access_token/token_type at top level.
	c.Header("AuthToken", resp.Auth.AccessToken)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login Successfully",
		"data":         resp,
		"success":      true,
		"access_token": resp.Auth.AccessToken,
		"token_type":   "bearer",
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RefreshRequest true "Refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(err.Error(), "VALIDATION_ERROR", nil))
		return
	}
	resp, err := h.authUseCase.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed(err.Error(), "UNAUTHORIZED", nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("", resp, nil))
}

// Logout godoc
// @Summary Logout current user
// @Tags auth
// @Security BearerAuth
// @Router /logout [get]
func (h *AuthHandler) Logout(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("token missing", "UNAUTHORIZED", nil))
		return
	}
	if err := h.authUseCase.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed(err.Error(), "INTERNAL_SERVER_ERROR", nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Successfully Logout", []any{}, nil))
}

// ChangePassword godoc
// @Summary Change current user's password
// @Tags auth
// @Security BearerAuth
// @Param body body dto.ChangePasswordRequest true "Old + new password"
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	principal, exists := middleware.Principal(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed("unauthorized", "UNAUTHORIZED", nil))
		return
	}
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(err.Error(), "VALIDATION_ERROR", nil))
		return
	}
	if err := h.authUseCase.ChangePassword(principal.User.ID, &req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(err.Error(), "BAD_REQUEST", nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Password changed successfully", []any{}, nil))
}

// UpdateProfile godoc
// @Summary Update current user's profile
// @Tags auth
// @Security BearerAuth
// @Param body body dto.ProfileUpdateRequest true "Profile fields"
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	principal, exists := middleware.Principal(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed("unauthorized", "UNAUTHORIZED", nil))
		return
	}
	var req dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(err.Error(), "VALIDATION_ERROR", nil))
		return
	}
	resp, err := h.authUseCase.UpdateProfile(principal.User.ID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(err.Error(), "BAD_REQUEST", nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Profile updated successfully", resp, nil))
}

// Register godoc
// @Summary Register new user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param register body dto.RegisterRequest true "Registration information"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUseCase.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("Register Success", resp, nil))
}

// GetUserPermissions godoc
// @Summary Get user permissions
// @Description Get permissions for the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PermissionSimple
// @Failure 401 {object} map[string]string
// @Router /auth/permissions [get]
func (h *AuthHandler) GetUserPermissions(c *gin.Context) {
	userID, exists := c.Get(constants.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID format"})
		return
	}
	permissions, err := h.authUseCase.GetUserPermissions(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.PermissionSimple
	for _, permission := range permissions {
		response = append(response, dto.PermissionSimple{
			ID:   permission.ID,
			Name: permission.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetUser
func (h *AuthHandler) GetUser(c *gin.Context) {

	userID, exists := c.Get(constants.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID format"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	token := ""
	if authHeader != "" && len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		token = authHeader[7:]
	}

	auth, err := h.authUseCase.GetUser(userUUID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Get User "+auth.User.FirstName, auth.User, nil))

}

func (h *AuthHandler) CreateMeta(c *gin.Context) {

	userID, exists := c.Get(constants.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID format"})
		return
	}

	var req dto.CreateMetaDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUseCase.CreateMetaData(userUUID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Create Meta Success", resp, nil))
}

func (h *AuthHandler) GetUserMeta(c *gin.Context) {
	userID, exists := c.Get(constants.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID format"})
		return
	}

	metaData, err := h.authUseCase.GetMetaData(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Get Meta Success", metaData, nil))
}

// func (h *AuthHandler) SendToMe(c *gin.Context) {
// 	var req dto.SendNotificationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get authenticated user
// 	user, exists := c.Get(constants.UserIDKey)
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
// 		return
// 	}

// 	// Send notification to authenticated user
// 	response, err := h.authUseCase.SendToUser(user.(*entities.User).ID, req.Title, req.Body, req.Data)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// func (h *AuthHandler) SendNotification(c *gin.Context) {
// 	var req dto.SendNotificationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var response *dto.NotificationResponse
// 	var err error

// 	// Send based on the provided parameters
// 	if req.Topic != "" {
// 		// Send to topic
// 		response, err = h.authUseCase.SendToTopic(req.Topic, req.Title, req.Body, req.Data)
// 	} else if len(req.UserIDs) > 0 {
// 		// Send to specific users
// 		response, err = h.authUseCase.SendToUsers(req.UserIDs, req.Title, req.Body, req.Data)
// 	} else {
// 		// Send to all
// 		response, err = h.authUseCase.SendToAll(req.Title, req.Body, req.Data)
// 	}

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// GetMyTokenHistory godoc
// @Summary Get current user's token history
// @Tags auth
// @Security BearerAuth
// @Router /auth/token-history [get]
func (h *AuthHandler) GetMyTokenHistory(c *gin.Context) {
	principal, exists := middleware.Principal(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed("unauthorized", "UNAUTHORIZED", nil))
		return
	}
	histories, err := h.authUseCase.GetTokenHistory(principal.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed(err.Error(), "INTERNAL_SERVER_ERROR", nil))
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("", histories, nil))
}
