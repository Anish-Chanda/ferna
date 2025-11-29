package auth

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateUserID generates a random UUID for user ID
func GenerateUserID() (uuid.UUID, error) {
	id := uuid.New()
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("failed to generate UUID")
	}
	return id, nil
}
