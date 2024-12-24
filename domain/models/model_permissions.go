package models

type ModelPermission struct {
	ID           uint   `gorm:"primaryKey"`
	ModelType    string `gorm:"size:100;not null"`
	ModelID      uint   `gorm:"not null"`
	PermissionID uint   `gorm:"not null"`
}
