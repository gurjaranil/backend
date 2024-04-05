package Entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string `valid:"required,min=3,max=50" json:"name" gorm:"not null"`
	Email         string `valid:"required,email" json:"email" gorm:"uniqueIndex;not null"`
	ContactNumber int    `valid:"required,numeric,length(10)" json:"contact_number" gorm:"not null"`
	Role          string `valid:"required" json:"role" gorm:"not null"`
	LibID         int    `valid:"required" json:"lib_id" gorm:"not null"`
	Password      string `valid:"required,min=6,max=20" json:"-" gorm:"not null"`

	Library Library `valid:"required" gorm:"foreignKey:lib_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"library"`
}
