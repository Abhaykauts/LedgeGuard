package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Abhaykauts/LedgeGuard/backend/internal/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// InitSQLite initializes a new SQLite connection
func InitSQLite(dbPath string) (*gorm.DB, error) {
	// Ensure the parent directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v", dir, err)
		} else {
			log.Printf("Successfully ensured directory exists: %s", dir)
		}
	}

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
