package service

import (
	"library/Entity"
	"library/utill"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type bookWithAvailablity struct {
	Book         Entity.BookInventory
	Availability string
}

// error in searching
func SearchBook(c *gin.Context) {
	keyword := c.Param("keyword")
	var books []Entity.BookInventory
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	result := db.Preload(clause.Associations).Where("isbn LIKE ? OR title LIKE ? OR authors LIKE ? OR publisher LIKE ?", keyword, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Where("lib_id = ?", lib).Find(&books)
	var booksWithAvailablity []bookWithAvailablity

	for i := range books {
		if books[i].AvailableCopies > 0 {
			var bookDetail bookWithAvailablity
			bookDetail.Book = books[i]
			bookDetail.Availability = "Available"
			booksWithAvailablity = append(booksWithAvailablity, bookDetail)
		} else {
			var bookDetail bookWithAvailablity
			bookDetail.Book = books[i]

			var issue Entity.IssueRegistery
			result := db.Where("reader_id IN (SELECT id from users where lib_id = ?)", lib).Where("expected_return_date > ?", time.Now()).Where("return_approver_id = ?", 0).First(&issue, "ISBN = ?", books[i].ISBN)
			if result.RowsAffected == 0 {
				bookDetail.Availability = "Can't say about availablity"

			} else {
				bookDetail.Availability = issue.ExpectedReturnDate.String()

			}

			booksWithAvailablity = append(booksWithAvailablity, bookDetail)
		}
	}
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found with for this keyword "+keyword)
		return
	} else {
		c.JSON(http.StatusOK, booksWithAvailablity)
	}

}
func CreateIssue(c *gin.Context) {
	var bookRequest struct {
		gorm.Model
		BookId int `valid:"required" json:"book_id"`
	}
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	if err := c.BindJSON(&bookRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {
		if !(utill.IsValidIsbn(bookRequest.BookId)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Enter valid isbn")
			return
		}

		var book Entity.BookInventory
		result := db.Preload(clause.Associations).Where("lib_id = ?", user.LibID).Find(&book, "isbn = ?", bookRequest.BookId)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
			return
		} else {

			var IssueRequest Entity.RequestEvents
			IssueRequest.BookID = bookRequest.BookId
			IssueRequest.ReaderID = int(user.ID)
			IssueRequest.RequestDate = time.Now()
			IssueRequest.Book = book
			IssueRequest.Reader = user

			if IssueRequest.Book.AvailableCopies > 0 {
				IssueRequest.RequestType = "issue"

				IssueRequest.RequestStatus = "pending"
			} else {
				IssueRequest.RequestType = "issue"

				IssueRequest.RequestStatus = "rejected"
			}
			tx := db.Begin()
			if err := tx.Create(&IssueRequest).Preload(clause.Associations).Error; err != nil {
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusNotFound, err)
				return
			} else {
				tx.Commit()
				if IssueRequest.RequestStatus == "rejected" {
					c.JSON(http.StatusNotAcceptable, "Request Rejected due to availablity")
					return
				} else {
					c.JSON(http.StatusCreated, IssueRequest)

				}

			}

		}
	}

}

func ReturnBook(c *gin.Context) {

	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	isbn := c.Param("isbn")
	var book Entity.BookInventory
	result := db.Preload(clause.Associations).Where("lib_id", lib).Find(&book, "ISBN = ?", isbn)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
		return
	} else {
		var issue Entity.IssueRegistery
		result2 := db.Preload(clause.Associations).Where("reader_id = ? AND return_approver_id = ?", user.ID, 0).Find(&issue, "ISBN = ?", isbn)
		if result2.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, "Book is not issued to you")
			return
		} else {
			var request Entity.RequestEvents
			request.RequestType = "return"
			request.RequestStatus = "pending"
			request.BookID = issue.BookId
			request.ReaderID = int(user.ID)
			request.RequestDate = time.Now()

			tx := db.Begin()

			if err := tx.Create(&request).Error; err != nil {
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			} else {
				tx.Commit()
				c.JSON(http.StatusCreated, "Return Request added")
			}
		}

	}
}
