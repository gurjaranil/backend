package database

import (
	"fmt"
	"library/Entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data/library.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	} else {
		fmt.Println("Connected to the database")
	}
	db.AutoMigrate(&Entity.User{})
	db.AutoMigrate(&Entity.Library{})
	db.AutoMigrate(&Entity.BookInventory{})
	db.AutoMigrate(&Entity.RequestEvents{})
	db.AutoMigrate(&Entity.IssueRegistery{})

	return db
}
