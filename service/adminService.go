package service

import (
	"crypto/tls"
	"library/Entity"
	"library/config"
	"library/utill"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetUsers(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	var users []Entity.User
	result := db.Preload("Library").Find(&users, "lib_id = ?", lib)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "No records found")
		return
	} else {
		c.JSON(http.StatusOK, users)

	}

}

func AddBook(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	library := user.Library

	var bookDetailes struct {
		ISBN        int
		Title       string
		Authors     string
		Publisher   string
		Version     string
		TotalCopies int
	}
	var book Entity.BookInventory
	if err := c.BindJSON(&bookDetailes); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {

		bookDetailes.Authors = strings.Trim(bookDetailes.Authors, " ")
		bookDetailes.Title = strings.Trim(bookDetailes.Title, " ")
		bookDetailes.Publisher = strings.Trim(bookDetailes.Publisher, " ")
		bookDetailes.Version = strings.Trim(bookDetailes.Version, " ")

		if !(utill.IsValidIsbn(bookDetailes.ISBN)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Isbn")
			return
		}
		if !(utill.IsNameValid(bookDetailes.Title)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Title")
			return
		}
		if !(utill.IsNameValid(bookDetailes.Authors)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Author name")
			return
		}
		if !(utill.IsNameValid(bookDetailes.Publisher)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Publisher name")
			return
		}
		if !(utill.IsNameValid(bookDetailes.Version)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Version")
			return
		}
		var ExistingBook Entity.BookInventory
		result := db.Preload(clause.Associations).Where("lib_id = ? ", lib).Find(&ExistingBook, "ISBN = ?", bookDetailes.ISBN)
		if result.RowsAffected != 0 {
			if ExistingBook.ISBN == bookDetailes.ISBN && ExistingBook.Version == bookDetailes.Version {
				ExistingBook.TotalCopies = ExistingBook.TotalCopies + bookDetailes.TotalCopies
				ExistingBook.AvailableCopies = ExistingBook.AvailableCopies + bookDetailes.TotalCopies
				ExistingBook.QrCode = utill.GenerateQRCode(ExistingBook)

				tx := db.Begin()

				if err := tx.Save(&ExistingBook).Error; err != nil {
					tx.Rollback()
					c.AbortWithStatusJSON(http.StatusInternalServerError, "Somthing went wrong")
					return
				} else {
					tx.Commit()
					c.JSON(http.StatusAccepted, ExistingBook)

				}

			}
		} else {
			book.Authors = bookDetailes.Authors
			book.ISBN = bookDetailes.ISBN
			book.LibID = lib
			book.Library = library
			book.Title = bookDetailes.Title
			book.Publisher = bookDetailes.Publisher
			book.Version = bookDetailes.Version
			book.TotalCopies = bookDetailes.TotalCopies
			book.AvailableCopies = book.TotalCopies
			book.QrCode = utill.GenerateQRCode(book)

			tx := db.Begin()

			if err := tx.Create(&book).Error; err != nil {
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusInternalServerError, "got some error")
				return
			} else {
				tx.Commit()
				c.JSON(http.StatusCreated, book)

			}

		}

	}
}
func RemoveBook(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	isbn := c.Param("ISBN")
	var ExistingBook Entity.BookInventory
	db.Where("lib_id", lib).Find(&ExistingBook, "isbn = ?", isbn)

	tx := db.Begin()
	if err := tx.Unscoped().Delete(&ExistingBook).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
		return
	} else {
		tx.Commit()
		c.JSON(http.StatusAccepted, ExistingBook.Title+" is removed")

	}
}
func BookById(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	isbn := c.Param("ISBN")
	var ExistingBook Entity.BookInventory
	result := db.Preload(clause.Associations).Where("lib_id", lib).Find(&ExistingBook, "ISBN = ?", isbn)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
		return
	} else {
		c.JSON(http.StatusOK, ExistingBook)

	}
}
func GetAllBooks(c *gin.Context) {
	var books []Entity.BookInventory
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	result := db.Preload("Library").Where("lib_id", lib).Find(&books)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "No book Available")
		return
	} else {
		c.JSON(http.StatusOK, books)

	}
}

func UpdateBook(c *gin.Context) {
	var bookDetailes struct {
		Title     string
		Authors   string
		Publisher string
		Version   string
	}

	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	isbn := c.Param("ISBN")

	var ExistingBook Entity.BookInventory
	if err := c.BindJSON(&bookDetailes); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {

		bookDetailes.Authors = strings.Trim(bookDetailes.Authors, " ")
		bookDetailes.Title = strings.Trim(bookDetailes.Title, " ")
		bookDetailes.Publisher = strings.Trim(bookDetailes.Publisher, " ")
		bookDetailes.Version = strings.Trim(bookDetailes.Version, " ")

		if !(utill.IsNameValid(bookDetailes.Title)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Title")
			return
		}
		if !(utill.IsNameValid(bookDetailes.Authors)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Author name")
			return
		}
		if !(utill.IsNameValid(bookDetailes.Publisher)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Publisher name")
			return
		}
		if !(utill.IsNameValid(bookDetailes.Version)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Version")
			return
		}

		result := db.Preload(clause.Associations).Where("lib_id", lib).Find(&ExistingBook, "ISBN = ?", isbn)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
			return
		} else {
			ExistingBook.Authors = bookDetailes.Authors
			ExistingBook.Title = bookDetailes.Title
			ExistingBook.Publisher = bookDetailes.Publisher
			ExistingBook.Version = bookDetailes.Version

			tx := db.Begin()

			if err := tx.Save(&ExistingBook).Error; err != nil {
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusInternalServerError, err)
				return
			} else {
				tx.Commit()
				c.JSON(http.StatusAccepted, ExistingBook)

			}

		}
	}
}
func IssueRequest(c *gin.Context) {
	var Requests []Entity.RequestEvents
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	result := db.Preload(clause.Associations).Where("reader_id IN (SELECT id from users where lib_id = ?)", lib).Find(&Requests)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "No request found")
		return
	} else {

		c.JSON(http.StatusOK, Requests)
	}
}

func ApproveRequest(c *gin.Context) {
	var Request Entity.RequestEvents
	id := c.Param("rid")
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	var issue Entity.IssueRegistery

	result := db.Preload(clause.Associations).Where("reader_id IN (SELECT id from users where lib_id = ?)", lib).Find(&Request, "id = ?", id)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	} else {

		if Request.RequestStatus == "approved" {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, "Request already Accepted")
			return

		} else if Request.RequestStatus == "rejected" {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, "Request already Rejected")
			return
		}

		if Request.RequestType == "return" {

			result2 := db.Preload(clause.Associations).Where("reader_id", Request.Reader.ID).Find(&issue, "ISBN = ?", Request.Book.ISBN)
			if result2.RowsAffected == 0 {
				c.AbortWithStatusJSON(http.StatusNotFound, "Book is not issued to you")
				return
			} else {

				issue.ExpectedReturnDate = time.Now()
				issue.ReturnApproverID = int(user.ID)
				issue.ReturnDate = time.Now()
				if Request.Book.AvailableCopies+1 > Request.Book.TotalCopies {
					c.AbortWithStatusJSON(501, "Somthing went wrong need to  check the database")
					return
				}
				Request.Book.AvailableCopies = Request.Book.AvailableCopies + 1

			}
		}

		Request.Approver = user
		Request.ApproverID = int(user.ID)
		Request.RequestStatus = "approved"
		Request.ApprovalDate = time.Now()

		if Request.RequestType == "issue" {
			if Request.Book.AvailableCopies > 0 {
				Request.Book.AvailableCopies = Request.Book.AvailableCopies - 1
				issue.ISBN = Request.Book.ISBN
				issue.BookId = Request.BookID
				issue.ReaderID = Request.ReaderID
				issue.IssueApproverID = int(user.ID)
				issue.IssueStatus = Request.RequestStatus
				issue.IssueDate = time.Now()
				issue.ExpectedReturnDate = time.Now().Add(time.Hour * 24 * 14)
				issue.IssueApprover = user
				issue.Reader = Request.Reader
			} else {
				c.AbortWithStatusJSON(http.StatusNotFound, "Book Current not available")
				return
			}
		}

		tx := db.Begin()

		if err := tx.Save(&Request).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, Request)
			return
		} else if err := tx.Save(&issue).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"issue": issue, "req": Request})
			return
		} else if err := tx.Save(&Request.Book).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")

		} else {
			tx.Commit()
			c.JSON(http.StatusAccepted, "Successfully approved request")

		}

	}

}
func RejectRequest(c *gin.Context) {
	var Request Entity.RequestEvents
	id := c.Param("rid")
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	result := db.Preload("Library").Where("reader_id IN (SELECT id from users where lib_id = ?)", lib).Find(&Request, "id = ?", id)
	if Request.RequestStatus == "approved" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, "Request already Accepted")
		return
	} else if Request.RequestStatus == "rejected" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, "Request already Rejected")
		return
	}
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	} else {
		Request.Approver = user
		Request.ApproverID = int(user.ID)
		Request.RequestStatus = "rejected"
		Request.ApprovalDate = time.Now()

		tx := db.Begin()

		if err := tx.Save(&Request).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return

		} else {
			tx.Commit()
			c.JSON(http.StatusAccepted, "Successfully Rejected request")

		}

	}
}
func IssuedBook(c *gin.Context) {
	var books []Entity.IssueRegistery
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	result := db.Preload(clause.Associations).Where("reader_id IN (SELECT id from users where lib_id = ?)", lib).Find(&books, "return_approver_id = ?", 0)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Records not found")
	} else {
		c.JSON(http.StatusOK, books)

	}
}
func UserIssue(c *gin.Context) {
	var books []Entity.IssueRegistery

	var user2 Entity.User

	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	id := user.ID
	uid := c.Param("id")
	results := db.Preload(clause.Associations).Where("lib_id = ?", lib).Find(&user2, "id = ?", id)
	if results.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "User not found")
		return
	}
	result := db.Preload(clause.Associations).Find(&books, "reader_id = ?", uid)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Records not found")
		return
	} else {
		c.JSON(http.StatusOK, books)

	}
}
func DecreaseBook(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	isbn := c.Param("isbn")
	var ExistingBook Entity.BookInventory
	result := db.Where("lib_id", lib).Find(&ExistingBook, "isbn = ?", isbn)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
		return
	} else {
		if ExistingBook.AvailableCopies > 0 {
			ExistingBook.AvailableCopies = ExistingBook.AvailableCopies - 1
			ExistingBook.TotalCopies = ExistingBook.TotalCopies - 1

		} else {
			c.AbortWithStatusJSON(http.StatusNotFound, "Available copies not found")
			return
		}

	}
	tx := db.Begin()
	if err := tx.Save(&ExistingBook).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
		return
	} else {
		tx.Commit()
		c.JSON(http.StatusAccepted, "Book Quantity  Decreased Successfully!")

	}
}

func IncreaseBook(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID
	isbn := c.Param("isbn")
	var ExistingBook Entity.BookInventory
	result := db.Where("lib_id", lib).Find(&ExistingBook, "isbn = ?", isbn)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
		return
	} else {
		ExistingBook.AvailableCopies = ExistingBook.AvailableCopies + 1
		ExistingBook.TotalCopies = ExistingBook.TotalCopies + 1
	}

	tx := db.Begin()
	if err := tx.Save(&ExistingBook).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotFound, "Book not found")
		return
	} else {
		tx.Commit()
		c.JSON(http.StatusAccepted, "Book Quantity Increased Successfully!")

	}
}
func Invite(c *gin.Context) {
	var inviteDetails struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := c.BindJSON(&inviteDetails); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {

		userDetails, _ := c.Get("user")
		user1 := userDetails.(Entity.User)
		inviteDetails.Email = strings.ToLower(inviteDetails.Email)
		inviteDetails.Email = strings.Trim(inviteDetails.Email, " ")
		inviteDetails.Name = strings.Trim(inviteDetails.Name, " ")

		if !utill.IsValidEmail(inviteDetails.Email) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Email Format")
			return
		}
		if !utill.IsNameValid(inviteDetails.Name) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Name")
			return
		}

		var user Entity.User
		result := db.Preload(clause.Associations).First(&user, "email = ?", inviteDetails.Email)
		if result.RowsAffected != 0 {
			c.JSON(http.StatusConflict, "User already a "+user.Role+" !")
			return
		}

		name := c.PostForm("name")
		m := gomail.NewMessage()
		m.SetHeader("From", os.Getenv("EMAIL"))
		m.SetHeader("To", inviteDetails.Email)
		m.SetHeader("Subject", "Invitation to join Library Management System")
		lib := strconv.Itoa(user1.LibID)
		m.SetBody("text/html", utill.GenerateInvitationEmailHTML(name, config.AppConfig.BaseURL+":"+config.AppConfig.Port+"/signupForm?email="+inviteDetails.Email+"&name="+inviteDetails.Name+"&library="+lib, user1.Email, user1.Library.Name))
		d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		if err := d.DialAndSend(m); err != nil {
			c.JSON(http.StatusInternalServerError, "Invitation Failed !")
			return
		} else {
			c.JSON(http.StatusOK, "Invitation sent succesfully !")

		}

	}

}
