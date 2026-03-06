package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}

	if hash == password {
		t.Fatal("HashPassword() returned unhashed password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecretPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Test correct password
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}
	if !match {
		t.Fatal("CheckPasswordHash() should return true for correct password")
	}

	// Test wrong password
	match, err = CheckPasswordHash("wrongPassword", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}
	if match {
		t.Fatal("CheckPasswordHash() should return false for wrong password")
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "samePassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Same password should produce different hashes (due to random salt)
	if hash1 == hash2 {
		t.Fatal("HashPassword() should produce unique hashes for same password")
	}

	// But both should still match the original password
	match1, _ := CheckPasswordHash(password, hash1)
	match2, _ := CheckPasswordHash(password, hash2)

	if !match1 || !match2 {
		t.Fatal("Both hashes should match the original password")
	}
}

func TestEmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	match, err := CheckPasswordHash("", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}
	if !match {
		t.Fatal("Empty password should match its hash")
	}

	match, _ = CheckPasswordHash("notEmpty", hash)
	if match {
		t.Fatal("Non-empty password should not match empty password hash")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	if token == "" {
		t.Fatal("MakeJWT() returned empty token")
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}

	if parsedID != userID {
		t.Fatalf("ValidateJWT() returned wrong userID: got %v, want %v", parsedID, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	// Create token that expires immediately (negative duration)
	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("ValidateJWT() should return error for expired token")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("ValidateJWT() should return error for wrong secret")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := ValidateJWT("invalid-token-string", "secret")
	if err == nil {
		t.Fatal("ValidateJWT() should return error for invalid token")
	}
}

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		name        string
		headers     http.Header
		wantToken   string
		expectError bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid-token-123"},
			},
			wantToken:   "valid-token-123",
			expectError: false,
		},
		{
			name:        "Missing Authorization header",
			headers:     http.Header{},
			wantToken:   "",
			expectError: true,
		},
		{
			name: "Invalid Authorization header format",
			headers: http.Header{
				"Authorization": []string{"InvalidFormat token"},
			},
			wantToken:   "",
			expectError: true,
		},
		{
			name: "Bearer prefix but no token",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			wantToken:   "",
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GetBearerToken(tc.headers)
			if tc.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if token != tc.wantToken {
				t.Fatalf("GetBearerToken() = %v, want %v", token, tc.wantToken)
			}
		})
	}

}
