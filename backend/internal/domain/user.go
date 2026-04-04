package domain

import "time"

// Role represents the user access level
type Role string

const (
	RoleViewer  Role = "VIEWER"
	RoleAnalyst Role = "ANALYST"
	RoleAdmin   Role = "ADMIN"
)

// User represents a system user
type User struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRepository defines the persistence contract
type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List() ([]User, error)
}
