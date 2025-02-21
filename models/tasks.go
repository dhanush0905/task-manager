package models

import "gorm.io/gorm"

// Task model represents a task entity
type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"` // Associates task with a user
}
