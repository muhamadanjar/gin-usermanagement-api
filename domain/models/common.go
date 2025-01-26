package models

import (
	"time"

	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Timestamp struct {
	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type (
	PaginationRequest struct {
		Search  string `form:"search"`
		Page    int    `form:"page"`
		PerPage int    `form:"per_page"`
	}

	PaginationResponse struct {
		Page    int   `json:"page"`
		PerPage int   `json:"per_page"`
		MaxPage int64 `json:"max_page"`
		Count   int64 `json:"count"`
	}
)

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *PaginationResponse) GetLimit() int {
	return p.PerPage
}

func (p *PaginationResponse) GetPage() int {
	return p.Page
}
