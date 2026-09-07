package entities

import (
	"time"

	"github.com/google/uuid"
)

type ModelPermission struct {
	ID           uuid.UUID
	ModelID      uuid.UUID // may be a "role", "menu", or another model type
	ModelType    string
	PermissionID uuid.UUID
	Permission   Permission
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
