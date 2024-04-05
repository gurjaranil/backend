package Entity

import "gorm.io/gorm"

type BookInventory struct {
	gorm.Model
	ISBN            int    `valid:"required" gorm:"primaryKey;not null" json:"isbn"`
	LibID           int    `valid:"required" gorm:"primaryKey;not null" json:"lib_id"`
	Title           string `valid:"required" json:"title" gorm:"not null"`
	Authors         string `valid:"required" json:"authors" gorm:"not null"`
	Publisher       string `valid:"required" json:"publisher" gorm:"not null"`
	Version         string `valid:"required" json:"version" gorm:"not null"`
	TotalCopies     int    `valid:"required" json:"total_copies" gorm:"not null"`
	AvailableCopies int    `valid:"required" json:"available_copies" gorm:"not null"`
	QrCode          string `valid:"required" json:"Qr_code" gorm:"not null"`

	Library Library `valid:"required" gorm:"foreignKey:lib_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"library"`
}
