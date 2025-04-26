// internal/delivery/http/handlers/notification_handler.go
package handlers

import (
	"net/http"
	"usermanagement-api/domain/entities"
	"usermanagement-api/internal/constants"
	"usermanagement-api/internal/dto"
	"usermanagement-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationUseCase usecase.NotificationUseCase
}

func NewNotificationHandler(notificationUseCase usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

// SendNotification godoc
// @Summary Send notification
// @Description Send a push notification to users
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body dto.SendNotificationRequest true "Notification details"
// @Success 200 {object} dto.NotificationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /notifications/send [post]
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response *dto.NotificationResponse
	var err error

	// Send based on the provided parameters
	if req.Topic != "" {
		// Send to topic
		response, err = h.notificationUseCase.SendToTopic(req.Topic, req.Title, req.Body, req.Data)
	} else if len(req.UserIDs) > 0 {
		// Send to specific users
		response, err = h.notificationUseCase.SendToUsers(req.UserIDs, req.Title, req.Body, req.Data)
	} else {
		// Send to all
		response, err = h.notificationUseCase.SendToAll(req.Title, req.Body, req.Data)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SendToMe godoc
// @Summary Send notification to self
// @Description Send a push notification to the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body dto.SendNotificationRequest true "Notification details"
// @Success 200 {object} dto.NotificationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /notifications/send-to-me [post]
func (h *NotificationHandler) SendToMe(c *gin.Context) {
	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user
	user, exists := c.Get(constants.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	// Send notification to authenticated user
	response, err := h.notificationUseCase.SendToUser(user.(*entities.User).ID, req.Title, req.Body, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
