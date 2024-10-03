package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Token struct {
	ID           int64      `json:"id" db:"id"`                           // Primary key
	UserID       int64      `json:"user_id" db:"user_id"`                 // Foreign key to the user table
	RefreshToken string     `json:"refresh_token" db:"refresh_token"`     // Refresh token string
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`           // Timestamp of creation
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`           // Timestamp of last update
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"` // Nullable timestamp for soft deletion
}
