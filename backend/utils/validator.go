// utils/validator.go
// A utility for handling request payload validation.
package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

func ParseAndValidate(r *http.Request, payload interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return fmt.Errorf("invalid request payload: %w", err)
	}
	if err := validate.Struct(payload); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
