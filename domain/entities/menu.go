package entities

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	ID            uuid.UUID
	Name          string
	Url           string
	PermissionKey string
	Icon          string
	Description   string
	ParentID      *uuid.UUID
	Parent        *Menu
	Children      []*Menu
	Permissions   []*Permission
	Sequence      int
	IsActive      bool
	IsVisible     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
