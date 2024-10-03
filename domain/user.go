package domain

import (
	"time"
)

type User struct {
	ID              int64      `json:"id"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"` // Nullable timestamp for soft deletion
	EmailVerifiedAt *time.Time `json:"email_verified_at"`                    // Change to pointer
}
