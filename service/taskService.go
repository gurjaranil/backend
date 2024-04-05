package service

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// service connection with database
func Connect(data *gorm.DB) {
	db = data
}
