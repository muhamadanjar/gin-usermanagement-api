package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;unique;not null"`
	Password string `gorm:"not null"`
	RoleID   uint   `gorm:"not null"`
}
