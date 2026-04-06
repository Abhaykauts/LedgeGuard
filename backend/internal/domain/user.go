package domain

import (
	"time"

	"gorm.io/gorm"
)

// Role represents the user access level
type Role string

const (
	RoleViewer  Role = "VIEWER"
	RoleAnalyst Role = "ANALYST"
	RoleAdmin   Role = "ADMIN"
)

// User represents a system user
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username" binding:"required,alphanum,min=3,max=20"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         Role           `gorm:"not null" json:"role" binding:"required,oneof=VIEWER ANALYST ADMIN"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
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
