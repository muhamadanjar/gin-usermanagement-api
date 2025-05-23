package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"usermanagement-api/internal/constants"
	"usermanagement-api/internal/dto"
	"usermanagement-api/internal/usecase"
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUseCase.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp, "message": "Login Success"})
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

	c.JSON(http.StatusCreated, resp)
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
	fmt.Println("user id", userID)
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

// CreateModelPermission godoc
// @Summary Create model permission
// @Description Assign a permission to a model (role or menu)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param model-permission body dto.ModelPermissionRequest true "Model permission information"
// @Success 201 {object} dto.ModelPermissionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/model-permissions [post]
func (h *AuthHandler) CreateModelPermission(c *gin.Context) {
	var req dto.ModelPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUseCase.CreateModelPermission(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetModelPermissions godoc
// @Summary Get model permissions
// @Description Get permissions for a specific model (role or menu)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param model_type query string true "Model type (role or menu)"
// @Param model_id query int true "Model ID"
// @Success 200 {array} dto.ModelPermissionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/model-permissions [get]
func (h *AuthHandler) GetModelPermissions(c *gin.Context) {
	modelType := c.Query("model_type")
	modelIDStr := c.Query("model_id")

	if modelType == "" || modelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "model_type and model_id are required"})
		return
	}

	// Removed unused modelID parsing

	modelUUID, err := uuid.Parse(modelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid model_id format"})
		return
	}
	permissions, err := h.authUseCase.GetModelPermissions(modelType, modelUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
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
	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Get User "+auth.User.FirstName, auth.User))

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
	c.JSON(http.StatusOK, gin.H{"data": resp, "message": "Create Meta Success"})
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
