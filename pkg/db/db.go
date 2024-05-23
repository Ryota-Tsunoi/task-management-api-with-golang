package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
