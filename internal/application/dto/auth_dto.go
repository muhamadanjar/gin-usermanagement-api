package dto

import "github.com/google/uuid"

type LoginRequest struct {
	// Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthInfoResponse struct {
	Auth AuthResponse `json:"auth"`
	User UserResponse `json:"user"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ModelPermissionRequest struct {
	ModelID      uuid.UUID `json:"model_id" binding:"required"`
	ModelType    string    `json:"model_type" binding:"required"`
	PermissionID uuid.UUID `json:"permission_id" binding:"required"`
}

type ModelPermissionResponse struct {
	ID           uuid.UUID        `json:"id"`
	ModelID      uuid.UUID        `json:"model_id"`
	ModelType    string           `json:"model_type"`
	PermissionID uuid.UUID        `json:"permission_id"`
	Permission   PermissionSimple `json:"permission"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

type CreateMetaDataRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type SendNotificationRequest struct {
	Title   string            `json:"title" binding:"required"`
	Body    string            `json:"body" binding:"required"`
	Data    map[string]string `json:"data"`
	UserIDs []uuid.UUID       `json:"user_ids,omitempty"`
	Topic   string            `json:"topic,omitempty"`
}

type NotificationResponse struct {
	Success      bool   `json:"success"`
	SuccessCount int    `json:"success_count,omitempty"`
	MessageID    string `json:"message_id,omitempty"`
	Error        string `json:"error,omitempty"`
}
