package model

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider represents the authentication provider enum from the SQL schema
type AuthProvider string

const (
	AuthProviderLocal    AuthProvider = "local"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderFacebook AuthProvider = "facebook"
)

// User represents the user model matching the SQL table schema
type User struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	AvatarURL    *string      `json:"avatar_url" db:"avatar_url"`
	AuthProvider AuthProvider `json:"auth_provider" db:"auth_provider"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	Email        string       `json:"email" db:"email"`
	FullName     string       `json:"full_name" db:"full_name"`
	PasswordHash *string      `json:"password_hash,omitempty" db:"password_hash"`
	Timezone     string       `json:"timezone" db:"timezone"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}