package domain

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100;unique;not null"`
}
