// internal/delivery/http/handlers/user_meta_handler.go
package handlers

import (
	"net/http"
	"usermanagement-api/internal/dto"
	"usermanagement-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserMetaHandler struct {
	userMetaUseCase usecase.UserMetaUseCase
}

func NewUserMetaHandler(userMetaUseCase usecase.UserMetaUseCase) *UserMetaHandler {
	return &UserMetaHandler{
		userMetaUseCase: userMetaUseCase,
	}
}

func (h *UserMetaHandler) CreateOrUpdate(c *gin.Context) {
	var req dto.CreateUserMetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userMetaUseCase.CreateOrUpdate(req.UserID, req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User meta updated successfully"})
}

func (h *UserMetaHandler) GetByKey(c *gin.Context) {
	userIDStr := c.Param("user_id")
	key := c.Param("key")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	resp, err := h.userMetaUseCase.GetByKey(userID, key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user meta not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserMetaHandler) GetAllByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	metaMap, err := h.userMetaUseCase.GetAllByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metaMap)
}

func (h *UserMetaHandler) Delete(c *gin.Context) {
	userIDStr := c.Param("user_id")
	key := c.Param("key")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userMetaUseCase.Delete(userID, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
