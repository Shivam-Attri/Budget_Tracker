// handlers/users.go
// Contains handlers for user registration, login, token refresh, and logout.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"your_username/budget-tracker/auth"
	"your_username/budget-tracker/database"
	"your_username/budget-tracker/models"
	"your_username/budget-tracker/utils"
)

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (env *Env) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload
	if err := utils.ParseAndValidate(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := &models.User{Email: payload.Email, Password: payload.Password}
	if err := env.DB.CreateUser(user); err != nil {
		if errors.Is(err, database.ErrEmailExists) {
			http.Error(w, "Email already in use", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (env *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, err := env.DB.GetUserByEmail(payload.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	accessToken, refreshToken, err := auth.GenerateTokens(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}
	refreshClaims, _ := auth.ValidateToken(refreshToken)
	// FIXED: Access ExpiresAt through the embedded RegisteredClaims struct.
	err = env.DB.StoreRefreshToken(user.ID, refreshToken, refreshClaims.RegisteredClaims.ExpiresAt.Time)
	if err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (env *Env) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	claims, err := auth.ValidateToken(payload.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	isValid, err := env.DB.ValidateRefreshToken(claims.UserID, payload.RefreshToken)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isValid {
		http.Error(w, "Invalid or revoked refresh token", http.StatusUnauthorized)
		return
	}
	if err := env.DB.RevokeRefreshToken(payload.RefreshToken); err != nil {
		http.Error(w, "Failed to revoke old token", http.StatusInternalServerError)
		return
	}
	newAccessToken, newRefreshToken, err := auth.GenerateTokens(claims.UserID)
	if err != nil {
		http.Error(w, "Failed to generate new tokens", http.StatusInternalServerError)
		return
	}
	newRefreshClaims, _ := auth.ValidateToken(newRefreshToken)
	// FIXED: Access ExpiresAt through the embedded RegisteredClaims struct.
	err = env.DB.StoreRefreshToken(claims.UserID, newRefreshToken, newRefreshClaims.RegisteredClaims.ExpiresAt.Time)
	if err != nil {
		http.Error(w, "Failed to store new refresh token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (env *Env) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := env.DB.RevokeRefreshToken(payload.RefreshToken); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}
