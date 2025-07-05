// models/category.go
// Defines the data structure for a user-created spending category.
package models

import "time"

type Category struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"` // Foreign key to the users table, hidden from JSON output
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
