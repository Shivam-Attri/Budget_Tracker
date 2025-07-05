// models/filter.go
// Defines a struct to hold all possible filter and sort parameters for transactions.
package models

import "time"

type TransactionFilter struct {
	Type      string
	StartDate *time.Time
	EndDate   *time.Time
	MinAmount *float64
	MaxAmount *float64
	Sort      string
}
