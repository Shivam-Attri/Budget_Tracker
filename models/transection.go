// models/transaction.go
// Defines the data structure for a single financial transaction.
package models

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"-"` // Foreign key to the users table, hidden from JSON output
	Description string          `json:"description" validate:"required"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	Date        time.Time       `json:"date" validate:"required"`
	Type        TransactionType `json:"type" validate:"required,oneof=income expense"`
	CategoryID  string          `json:"category_id" validate:"required,uuid"`
}
