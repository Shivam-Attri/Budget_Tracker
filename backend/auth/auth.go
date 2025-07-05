// auth/auth.go
// Handles JWT generation/validation, password hashing, and auth middleware.
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
)

// Init initializes the auth package with secrets and TTLs.
func Init(secret string, accessTTLMinutes, refreshTTLHours int) {
	jwtKey = []byte(secret)
	accessTokenTTL = time.Minute * time.Duration(accessTTLMinutes)
	refreshTokenTTL = time.Hour * time.Duration(refreshTTLHours)
}

// Claims defines the structure for access token claims.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateTokens creates a new pair of access and refresh tokens.
func GenerateTokens(userID string) (accessToken string, refreshToken string, err error) {
	// Generate Access Token
	accessExpirationTime := time.Now().Add(accessTokenTTL)
	accessClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshExpirationTime := time.Now().Add(refreshTokenTTL)
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtKey)
	return
}

// ValidateToken validates a JWT and returns its claims.
func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// HashToken creates a SHA256 hash of a token string.
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Middleware validates the access token for protected routes.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Access token has expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid access token", http.StatusUnauthorized)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIDFromContext(r *http.Request) (string, error) {
	userID := r.Context().Value("userID")
	if userID == nil {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID.(string), nil
}
