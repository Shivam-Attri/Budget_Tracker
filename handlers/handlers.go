// handlers/handlers.go
// Defines the Env struct for dependency injection into handlers.
package handlers

import "your_username/budget-tracker/database"

type Env struct {
	DB                database.Store
	DefaultPageLimit  int
	AllowedSortFields []string
}
