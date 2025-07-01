package storage

import (
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Transaction{},
		&model.User{},
	)
}
