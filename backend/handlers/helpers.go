package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-pkgz/auth/v2/token"
)

// extracts the user ID from the JWT claims in the request context.
func GetUserIDFromRequest(r *http.Request) (int64, error) {
	// First check if this is a test context with mocked user
	// TODO: There must be a better way to test....
	if userCtx := r.Context().Value("user"); userCtx != nil {
		if userMap, ok := userCtx.(map[string]interface{}); ok {
			if uid, ok := userMap["uid"].(string); ok {
				return strconv.ParseInt(uid, 10, 64)
			}
		}
	}

	// Production code path using JWT token
	user, err := token.GetUserInfo(r)
	if err != nil {
		return 0, errors.New("Failed to get user info")
	}
	userID, err := strconv.ParseInt(user.StrAttr("uid"), 10, 64)
	if err != nil {
		return 0, errors.New("Invalid user ID in token")
	}
	return userID, nil
}
