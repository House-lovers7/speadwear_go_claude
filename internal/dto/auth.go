package dto

import "time"

// LoginRequest represents login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SignupRequest represents signup request body
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserResponse `json:"user"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Picture   string    `json:"picture,omitempty"`
	Admin     bool      `json:"admin"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordUpdateRequest represents password update request
type PasswordUpdateRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}