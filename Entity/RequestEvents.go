package Entity

import (
	"time"

	"gorm.io/gorm"
)

type RequestEvents struct {
	gorm.Model
	BookID        int       `validate:"required" json:"book_id" gorm:"not null"`
	ReaderID      int       `valid:"required"   json:"reader_id" gorm:"not null"`
	RequestDate   time.Time `valid:"required,date" json:"request_date" gorm:"not null"`
	ApprovalDate  time.Time `valid:"date" json:"approval_date" gorm:"not null"`
	ApproverID    int       `valid:"required" json:"approver_id" gorm:"not null"`
	RequestStatus string    `valid:"required" json:"status" gorm:"not null"`
	RequestType   string    `valid:"required" json:"type" gorm:"not null"`

	Reader   User          `valid:"-" gorm:"foreignKey:reader_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"reader"`
	Book     BookInventory `valid:"-" gorm:"foreignKey:book_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"book"`
	Approver User          `valid:"-" gorm:"foreignKey:approver_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"approver"`
}
