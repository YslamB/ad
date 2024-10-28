package models

import "time"

type CreateAd struct {
	ID          int       `json:"id" binding:"omitempty"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	CreatedAt   time.Time `json:"created_at" binding:"omitempty"`
	IsActive    bool      `json:"is_active" binding:"omitempty"`
}
