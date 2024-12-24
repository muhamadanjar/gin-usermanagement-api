package models

type Menu struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100;not null"`
	URL          string `gorm:"size:200;not null"`
	PermissionID uint   `gorm:"not null"`
}
