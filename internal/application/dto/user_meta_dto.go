// internal/dto/user_meta_dto.go
package dto

import "github.com/google/uuid"

type CreateUserMetaRequest struct {
	Key    string    `json:"key" binding:"required"`
	Value  string    `json:"value" binding:"required"`
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type UpdateUserMetaRequest struct {
	Value string `json:"value" binding:"required"`
}

type UserMetaResponse struct {
	ID     uint      `json:"id"`
	Key    string    `json:"key"`
	Value  string    `json:"value"`
	UserID uuid.UUID `json:"user_id"`
}
