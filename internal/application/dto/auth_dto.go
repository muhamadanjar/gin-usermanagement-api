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

// RefreshRequest mirrors FastAPI's RefreshRequest {refresh_token}.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ProfileUpdateRequest struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"omitempty,email"`
	Avatar    string `json:"avatar"`
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
