package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := "24h"

	tests := []struct {
		name       string
		userID     uint
		email      string
		secret     string
		expiration string
		wantErr    bool
	}{
		{
			name:       "valid user",
			userID:     123,
			email:      "test@example.com",
			secret:     secret,
			expiration: expiration,
			wantErr:    false,
		},
		{
			name:       "zero user ID",
			userID:     0,
			email:      "zero@example.com",
			secret:     secret,
			expiration: expiration,
			wantErr:    false,
		},
		{
			name:       "large user ID",
			userID:     4294967295, // Max uint32
			email:      "large@example.com",
			secret:     secret,
			expiration: expiration,
			wantErr:    false,
		},
		{
			name:       "empty secret",
			userID:     123,
			email:      "test@example.com",
			secret:     "",
			expiration: expiration,
			wantErr:    false, // Empty secret still generates a token
		},
		{
			name:       "invalid expiration",
			userID:     123,
			email:      "test@example.com",
			secret:     secret,
			expiration: "invalid",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email, tt.secret, tt.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check that token is not empty
				if token == "" {
					t.Error("GenerateToken() returned empty token")
				}

				// Verify token structure (should have 3 parts separated by dots)
				// header.payload.signature
				parts := len(token)
				dots := 0
				for i := 0; i < parts; i++ {
					if token[i] == '.' {
						dots++
					}
				}
				if dots != 2 {
					t.Errorf("GenerateToken() returned invalid token structure, expected 2 dots, got %d", dots)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret-key"
	expiration := "24h"

	// Generate a valid token for testing
	validUserID := uint(123)
	validEmail := "test@example.com"
	validToken, err := GenerateToken(validUserID, validEmail, secret, expiration)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// Generate token with different secret
	invalidSecretToken, _ := GenerateToken(validUserID, validEmail, "different-secret", expiration)

	// Generate expired token
	expiredToken := generateExpiredToken(t, validUserID, validEmail)

	tests := []struct {
		name       string
		token      string
		secret     string
		wantUserID uint
		wantEmail  string
		wantErr    bool
		errString  string
	}{
		{
			name:       "valid token",
			token:      validToken,
			secret:     secret,
			wantUserID: validUserID,
			wantEmail:  validEmail,
			wantErr:    false,
		},
		{
			name:       "empty token",
			token:      "",
			secret:     secret,
			wantUserID: 0,
			wantEmail:  "",
			wantErr:    true,
			errString:  "token contains an invalid number of segments",
		},
		{
			name:       "invalid token format",
			token:      "invalid.token",
			secret:     secret,
			wantUserID: 0,
			wantEmail:  "",
			wantErr:    true,
			errString:  "token contains an invalid number of segments",
		},
		{
			name:       "token with wrong secret",
			token:      invalidSecretToken,
			secret:     secret,
			wantUserID: 0,
			wantEmail:  "",
			wantErr:    true,
			errString:  "signature is invalid",
		},
		{
			name:       "expired token",
			token:      expiredToken,
			secret:     secret,
			wantUserID: 0,
			wantEmail:  "",
			wantErr:    true,
			errString:  "token has invalid claims: token is expired",
		},
		{
			name:       "malformed token",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
			secret:     secret,
			wantUserID: 0,
			wantEmail:  "",
			wantErr:    true,
			errString:  "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errString != "" {
				if !containsString(err.Error(), tt.errString) {
					t.Errorf("ValidateToken() error = %v, want error containing %v", err, tt.errString)
				}
			}

			if !tt.wantErr && claims != nil {
				if claims.UserID != tt.wantUserID {
					t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, tt.wantUserID)
				}
				if claims.Email != tt.wantEmail {
					t.Errorf("ValidateToken() Email = %v, want %v", claims.Email, tt.wantEmail)
				}
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	secret := "test-secret-key"
	shortExpiration := "1s" // 1 second expiration for testing
	email := "test@example.com"

	userID := uint(123)
	token, err := GenerateToken(userID, email, secret, shortExpiration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Token should be valid immediately
	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Errorf("Token should be valid immediately after generation: %v", err)
	}
	if claims != nil && claims.UserID != userID {
		t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, userID)
	}

	// Wait for token to expire
	time.Sleep(1100 * time.Millisecond) // Wait 1.1 seconds

	// Token should now be expired
	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Error("Token should be expired after wait time")
	}
}

func TestJWTSecretHandling(t *testing.T) {
	email := "test@example.com"
	expiration := "24h"
	
	// Test with empty secret (should still work, but not secure)
	_, err := GenerateToken(123, email, "", expiration)
	if err != nil {
		t.Error("GenerateToken() failed with empty secret:", err)
	}

	// Test token validation with wrong secret
	token, _ := GenerateToken(123, email, "secret1", expiration)
	_, err = ValidateToken(token, "secret2")
	if err == nil {
		t.Error("ValidateToken() should fail with wrong secret")
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	secret := "benchmark-secret-key"
	userID := uint(123)
	email := "bench@example.com"
	expiration := "24h"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateToken(userID, email, secret, expiration)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	secret := "benchmark-secret-key"
	token, err := GenerateToken(123, "bench@example.com", secret, "24h")
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ValidateToken(token, secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Helper functions

func generateExpiredToken(t *testing.T, userID uint, email string) string {
	t.Helper()

	secret := "test-secret-key"

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	return tokenString
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsString(s[1:], substr) || len(substr) > 0 && s[0] == substr[0] && containsString(s[1:], substr[1:]))
}