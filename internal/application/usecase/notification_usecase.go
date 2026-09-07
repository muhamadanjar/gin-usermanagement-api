// internal/usecase/notification_usecase.go
package usecase

import (
	"usermanagement-api/domain/ports"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/application/dto"

	"github.com/google/uuid"
)

type NotificationUseCase interface {
	SendToUser(userID uuid.UUID, title, body string, data map[string]string) (*dto.NotificationResponse, error)
	SendToUsers(userIDs []uuid.UUID, title, body string, data map[string]string) (*dto.NotificationResponse, error)
	SendToTopic(topic, title, body string, data map[string]string) (*dto.NotificationResponse, error)
	SendToAll(title, body string, data map[string]string) (*dto.NotificationResponse, error)
}

type notificationUseCase struct {
	userMetaRepo repositories.UserMetaRepository
	notifier     ports.Notifier
}

func NewNotificationUseCase(
	userMetaRepo repositories.UserMetaRepository,
	notifier ports.Notifier,
) NotificationUseCase {
	return &notificationUseCase{
		userMetaRepo: userMetaRepo,
		notifier:     notifier,
	}
}

func (uc *notificationUseCase) SendToUser(userID uuid.UUID, title, body string, data map[string]string) (*dto.NotificationResponse, error) {
	// Get user device tokens
	deviceTokens, err := uc.userMetaRepo.FindByUserIDAndKey(userID, "fcm_token")
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Error:   "Failed to find user devices",
		}, err
	}

	var tokens = deviceTokens.Value

	// Send notification
	successCount, err := uc.notifier.SendToDevices([]string{tokens}, title, body, data)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &dto.NotificationResponse{
		Success:      true,
		SuccessCount: successCount,
	}, nil
}

func (uc *notificationUseCase) SendToUsers(userIDs []uuid.UUID, title, body string, data map[string]string) (*dto.NotificationResponse, error) {
	var allTokens []string

	// Get device tokens for each user
	for _, userID := range userIDs {
		deviceTokens, err := uc.userMetaRepo.FindByUserIDAndKey(userID, "fcm_token")
		if err != nil {
			continue
		}

		if deviceTokens != nil {
			allTokens = append(allTokens, deviceTokens.Value)
		}
	}

	if len(allTokens) == 0 {
		return &dto.NotificationResponse{
			Success: false,
			Error:   "No devices registered for users",
		}, nil
	}

	// Send notification
	successCount, err := uc.notifier.SendToDevices(allTokens, title, body, data)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &dto.NotificationResponse{
		Success:      true,
		SuccessCount: successCount,
	}, nil
}

func (uc *notificationUseCase) SendToTopic(topic, title, body string, data map[string]string) (*dto.NotificationResponse, error) {
	messageID, err := uc.notifier.SendToTopic(topic, title, body, data)
	if err != nil {
		return &dto.NotificationResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &dto.NotificationResponse{
		Success:   true,
		MessageID: messageID,
	}, nil
}

func (uc *notificationUseCase) SendToAll(title, body string, data map[string]string) (*dto.NotificationResponse, error) {
	// For sending to all, we use a topic that all devices are subscribed to
	// You would need to manage this subscription separately
	return uc.SendToTopic("all_users", title, body, data)
}
