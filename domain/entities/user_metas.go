package entities

import "github.com/google/uuid"

type UserMeta struct {
	ID     uint
	Key    string
	Value  string
	UserID uuid.UUID
}
