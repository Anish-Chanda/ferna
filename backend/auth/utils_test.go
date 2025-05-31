package auth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHashPassword_Format verifies that HashPassword returns a string in the expected Argon2id format.
func TestHashPassword_Format(t *testing.T) {
	password := "myS3cretP@ss!"
	encoded1, err := HashPassword(password)
	require.NoError(t, err, "HashPassword should not error")

	// The encoded string must start with "$argon2id$"
	require.True(t, strings.HasPrefix(encoded1, "$argon2id$"), "Encoded string must start with `$argon2id$`")

	fields := strings.Split(encoded1, "$")
	// Expected segments: ["", "argon2id", "v=<version>", "m=<mem>,t=<time>,p=<threads>", "<salt>", "<hash>"]
	require.Len(t, fields, 6, "Encoded string should split into 6 parts")

	assert.Equal(t, "argon2id", fields[1], "Second field should be 'argon2id'")
	assert.True(t, strings.HasPrefix(fields[2], "v="), "Third field should start with 'v='")
	assert.True(t, strings.HasPrefix(fields[3], "m="), "Fourth field should start with 'm='")

	// Splitting again to ensure salt and hash are non-empty Base64 strings:
	salt := fields[4]
	hash := fields[5]
	assert.NotEmpty(t, salt, "Salt portion must not be empty")
	assert.NotEmpty(t, hash, "Hash portion must not be empty")

	// Calling HashPassword again on the same password should produce a different encoded string
	encoded2, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEqual(t, encoded1, encoded2, "Two calls to HashPassword on the same password should yield different salts/hashes")
}

// TestHashPassword_Empty verifies that hashing an empty password still returns a valid format
// and that verifying an empty password against its generated hash succeeds.
func TestHashPassword_Empty(t *testing.T) {
	password := ""
	encoded, err := HashPassword(password)
	require.NoError(t, err, "HashPassword on empty string should not error")

	fields := strings.Split(encoded, "$")
	require.Len(t, fields, 6, "Encoded string should split into 6 parts even for empty password")
	assert.Equal(t, "argon2id", fields[1], "Second field should be 'argon2id'")

	ok, err := VerifyPassword(password, encoded)
	require.NoError(t, err, "VerifyPassword should not error on empty password and its hash")
	assert.True(t, ok, "VerifyPassword must return true for the correct (empty) password")
}

// TestVerifyPassword_SuccessAndFailure verifies that VerifyPassword returns true when using
// the correct password-hash pair and false when the password is wrong.
func TestVerifyPassword_SuccessAndFailure(t *testing.T) {
	password := "Test1234!"
	wrongPassword := "WrongPass"
	encoded, err := HashPassword(password)
	require.NoError(t, err)

	// Correct password → must verify successfully
	ok, err := VerifyPassword(password, encoded)
	require.NoError(t, err, "VerifyPassword should not error for correct password")
	assert.True(t, ok, "VerifyPassword must return true when password is correct")

	// Incorrect password → VerifyPassword should return (false, nil)
	ok, err = VerifyPassword(wrongPassword, encoded)
	require.NoError(t, err, "VerifyPassword should not error for wrong password")
	assert.False(t, ok, "VerifyPassword must return false when password is incorrect")
}

// TestVerifyPassword_InvalidFormat ensures that VerifyPassword rejects strings that are
// not in the correct Argon2id `$argon2id$...` format.
func TestVerifyPassword_InvalidFormat(t *testing.T) {
	invalidHashes := []string{
		"",                               // completely empty
		"not-a-hash",                     // no dollar signs
		"$argon2id$",                     // too few segments
		"$argon2id$v=19$m=65536,t=1,p=4", // missing salt/hash segments
		"$bcrypt$v=2a$...",               // wrong algorithm identifier
	}

	for _, encoded := range invalidHashes {
		ok, err := VerifyPassword("anyPassword", encoded)
		assert.False(t, ok, fmt.Sprintf("VerifyPassword should return ok=false for invalid format: %q", encoded))
		require.Error(t, err, fmt.Sprintf("VerifyPassword must error on invalid format: %q", encoded))
		assert.Contains(t, err.Error(), "invalid hash format", "Error should mention invalid hash format")
	}
}

// TestVerifyPassword_MismatchedParams tests that if the Argon2 parameters are altered (e.g. memory/time),
// VerifyPassword will derive a different key and therefore return false.
func TestVerifyPassword_MismatchedParams(t *testing.T) {
	// Hash the password normally
	password := "paramTest!"
	encodedNormal, err := HashPassword(password)
	require.NoError(t, err)

	// Now manually craft a second “valid‐looking” encoded string but tweak the parameter field:
	fields := strings.Split(encodedNormal, "$")
	require.Len(t, fields, 6, "Original encoded must split into 6 parts")

	// Change memory parameter from 65536 to 131072 (double)
	mParams := strings.Replace(fields[3], fmt.Sprintf("m=%d", argonMemory), fmt.Sprintf("m=%d", argonMemory*2), 1)

	// Reassemble an encoded string that looks valid but has a mismatched 'm' value.
	encodedMismatched := "$" + strings.Join([]string{
		fields[1], // "argon2id"
		fields[2], // e.g. "v=19"
		mParams,   // tweaked memory
		fields[4], // same salt (Base64)
		fields[5], // same hash (Base64)
	}, "$")

	// Since parameters changed, VerifyPassword should recompute a different hash → return false
	ok, err := VerifyPassword(password, encodedMismatched)
	require.NoError(t, err, "VerifyPassword should not error even if parameters are mismatched")
	assert.False(t, ok, "VerifyPassword must return false when Argon2 parameters differ")
}
