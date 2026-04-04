package database

import (
	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// InitSQLite initializes a new SQLite connection
func InitSQLite(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-Migrate the entities
	err = db.AutoMigrate(&domain.User{}, &domain.Record{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
