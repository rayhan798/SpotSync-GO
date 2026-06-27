package dto

import "time"

// UserSummary - Response data for user information
type UserSummary struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRes - Response for successful user registration
type RegisterRes struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRes - Response for successful user login with JWT token
type LoginRes struct {
	Token string      `json:"token"`
	User  UserSummary `json:"user"`
}
