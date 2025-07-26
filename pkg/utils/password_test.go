package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt accepts empty strings
		},
		{
			name:     "long password",
			password: "verylongpasswordthatexceedsnormallengthrequirementsbutshoulstillwork",
			wantErr:  false,
		},
		{
			name:     "special characters",
			password: "p@ssw0rd!#$%",
			wantErr:  false,
		},
		{
			name:     "unicode characters",
			password: "パスワード123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check that hash is not empty
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}

				// Check that hash is different from original password
				if hash == tt.password {
					t.Error("HashPassword() returned unhashed password")
				}

				// Check that the same password generates different hashes
				hash2, err := HashPassword(tt.password)
				if err != nil {
					t.Errorf("HashPassword() second call error = %v", err)
				}
				if hash == hash2 {
					t.Error("HashPassword() generated identical hashes for same password")
				}
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	// First, create some test hashes
	validPassword := "password123"
	validHash, err := HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to create test hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "valid password and hash",
			password: validPassword,
			hash:     validHash,
			want:     true,
		},
		{
			name:     "invalid password",
			password: "wrongpassword",
			hash:     validHash,
			want:     false,
		},
		{
			name:     "empty password with valid hash",
			password: "",
			hash:     validHash,
			want:     false,
		},
		{
			name:     "valid password with invalid hash",
			password: validPassword,
			hash:     "invalidhash",
			want:     false,
		},
		{
			name:     "empty password and empty hash",
			password: "",
			hash:     "",
			want:     false,
		},
		{
			name:     "case sensitive check",
			password: "Password123", // Different case
			hash:     validHash,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.password, tt.hash); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordHashingIntegration(t *testing.T) {
	passwords := []string{
		"short",
		"mediumlengthpassword",
		"verylongpasswordwithmanydifferentcharactersandnumbers1234567890",
		"!@#$%^&*()_+-=[]{}|;:,.<>?",
		"日本語パスワード",
		"Mixed123!@#日本語",
	}

	for _, password := range passwords {
		t.Run("password: "+password, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			// Verify the password
			if !CheckPassword(password, hash) {
				t.Error("CheckPassword() failed to verify correct password")
			}

			// Verify wrong password fails
			if CheckPassword(password+"wrong", hash) {
				t.Error("CheckPassword() verified incorrect password")
			}
		})
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !CheckPassword(password, hash) {
			b.Fatal("password check failed")
		}
	}
}