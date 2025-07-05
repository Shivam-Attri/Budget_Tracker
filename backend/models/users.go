// models/user.go
// Defines the data structure for a user.
package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // The '-' tag prevents the password from being marshalled into JSON
	CreatedAt time.Time `json:"created_at"`
}
