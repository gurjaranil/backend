package Entity

import "gorm.io/gorm"

type Library struct {
	gorm.Model
	Name string `valid:"required,min=3,max=50"  json:"name" gorm:"uniqueIndex;not null"`
}
