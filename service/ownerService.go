package service

import (
	"fmt"
	"library/Entity"
	"library/utill"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterLibrary(c *gin.Context) {
	var data struct {
		Name          string `valid:"required"  json:"name"`
		Email         string `valid:"email"      json:"email"`
		ContactNumber int    `valid:"length(10)"  json:"contact_number"`
		Password      string `valid:"length(6|20)"  json:"password"`
		LibName       string `valid:"required"    json:"lib_name"`
	}

	var library Entity.Library
	var user Entity.User
	if err := c.BindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else {
		fmt.Println(data)
		data.Email = strings.Trim(data.Email, " ")
		data.Email = strings.ToLower(data.Email)
		data.Name = strings.Trim(data.Name, " ")
		data.Password = strings.Trim(data.Password, " ")
		data.LibName = strings.Trim(data.LibName, " ")
		if !(utill.IsValidEmail(data.Email)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid email")
			return
		}
		if !(utill.IsContactNumberValid(data.ContactNumber)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Contact numbr")
			return
		}
		if !(utill.IsNameValid(data.Name)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Name")
			return
		}
		if !(utill.IsNameValid(data.LibName)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid library name")
			return
		}

		if !(utill.IsPasswordValid(data.Password)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "password must be between 8 to 15 characters included capital letters smaller letters symbols and numbers ")
			return
		}

		library.Name = data.LibName

		tx := db.Begin()
		if err := tx.Create(&library).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusConflict, "Library already exists")
			return
		}
		user.Name = data.Name
		user.ContactNumber = data.ContactNumber
		user.Email = data.Email
		user.Role = "owner"
		user.LibID = int(library.ID)
		user.Library = library
		bs, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user.Password = string(bs)

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusConflict, "User already exists")
			return
		}
		if err := tx.Commit().Error; err == nil {
			c.JSON(http.StatusCreated, gin.H{"user": user})
		} else {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, "Not added")

		}

	}

}

func AddAdmin(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)

	lib := user.LibID
	library := user.Library

	var adminDetails struct {
		Name          string `json:"name"`
		ContactNumber int    `json:"contact_number"`
		Email         string `valid:"email"      json:"email" `
		Password      string `valid:"length(6|20)" json:"password"`
	}

	var admin Entity.User

	if err := c.BindJSON(&adminDetails); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	} else {

		adminDetails.Email = strings.Trim(adminDetails.Email, " ")
		adminDetails.Email = strings.ToLower(adminDetails.Email)

		adminDetails.Password = strings.Trim(adminDetails.Password, " ")
		adminDetails.Name = strings.Trim(adminDetails.Name, " ")

		if !(utill.IsValidEmail(adminDetails.Email)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Email Format")
			return
		}
		if !(utill.IsContactNumberValid(adminDetails.ContactNumber)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Contact number format enter without +91 and must be 10 number only")
			return
		}
		if !(utill.IsNameValid(adminDetails.Name)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Name Format")
			return
		}
		if !(utill.IsPasswordValid(adminDetails.Password)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "password must be between 8 to 15 characters included capital letters smaller letters symbols and numbers ")
			return
		}

		var existingAdmin Entity.User
		result := db.Where("lib_id = ?", lib).First(&existingAdmin, "role = ?", "admin")
		if result.RowsAffected == 0 {
			admin.Name = adminDetails.Name
			admin.ContactNumber = adminDetails.ContactNumber
			admin.Email = adminDetails.Email
			admin.Role = "admin"
			admin.LibID = lib
			admin.Library = library

			bs, err := bcrypt.GenerateFromPassword([]byte(adminDetails.Password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)

			}
			admin.Password = string(bs)
			tx := db.Begin()
			if err := tx.Create(&admin).Error; err != nil {
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"msg": "user already exists"})
				return
			}
			if err := tx.Commit().Error; err == nil {
				c.JSON(http.StatusCreated, gin.H{"admin": admin})

			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"msg": "admin already exists"})
			return
		}
	}

}

func RemoveAdmin(c *gin.Context) {
	userDetails, _ := c.Get("user")
	user := userDetails.(Entity.User)
	lib := user.LibID

	var ExistingAdmin Entity.User
	db.Where("role", "admin").Find(&ExistingAdmin, "lib_id = ?", lib)
	tx := db.Begin()

	result := tx.Unscoped().Delete(&ExistingAdmin).Error
	if result != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotFound, "Admin not found")
		return
	} else {
		tx.Commit()
		c.JSON(http.StatusAccepted, ExistingAdmin.Name+" remove as admin")

	}

}
