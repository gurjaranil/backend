package Entity

import (
	"time"

	"gorm.io/gorm"
)

type IssueRegistery struct {
	gorm.Model
	ISBN               int       `valid:"required" json:"isbn" gorm:"not null"`
	BookId             int       `valid:"required" json:"book_id" gorm:"not null"`
	ReaderID           int       `valid:"required"   json:"reader_id" gorm:"not null"`
	IssueStatus        string    `valid:"required" json:"status" gorm:"not null"`
	IssueDate          time.Time `valid:"required,date" json:"issue_date" gorm:"not null"`
	ExpectedReturnDate time.Time `valid:"required,date" json:"expected_return_date" gorm:"not null"`
	ReturnDate         time.Time `valid:"date" json:"return_date" `
	ReturnApproverID   int       `json:"return_approver_id" `
	IssueApproverID    int       `valid:"required" json:"issue_approver_id" gorm:"not null"`

	Book           BookInventory `valid:"required" gorm:"foreignKey:book_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"book"`
	Reader         User          `valid:"-" gorm:"foreignKey:reader_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"reader"`
	ReturnApprover User          `valid:"-" gorm:"foreignKey:return_approver_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"return_approver"`
	IssueApprover  User          `valid:"-" gorm:"foreignKey:issue_approver_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"issue_approver"`
}
