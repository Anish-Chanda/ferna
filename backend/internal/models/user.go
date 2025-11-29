package models

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider represents the authentication provider
type AuthProvider string

const (
	AuthProviderLocal    AuthProvider = "local"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderFacebook AuthProvider = "facebook"
)

// User represents the user model matching the SQL schema
type User struct {
	ID           uuid.UUID    `json:"id"`
	AvatarURL    *string      `json:"avatar_url"`
	AuthProvider AuthProvider `json:"auth_provider"`
	CreatedAt    time.Time    `json:"created_at"`
	Email        string       `json:"email"`
	FullName     string       `json:"full_name"`
	PasswordHash *string      `json:"password_hash,omitempty"`
	Timezone     string       `json:"timezone"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
