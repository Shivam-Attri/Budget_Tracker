// models/budget.go
// Defines the data structure for a monthly budget.
package models

import "time"

type Budget struct {
	ID         string    `json:"id"`
	UserID     string    `json:"-"` // Foreign key to the users table, hidden from JSON output
	CategoryID string    `json:"category_id" validate:"required,uuid"`
	Amount     float64   `json:"amount" validate:"required,gt=0"`
	Month      int       `json:"month" validate:"required,min=1,max=12"`
	Year       int       `json:"year" validate:"required,min=2000"`
	CreatedAt  time.Time `json:"created_at"`
}
