package service

import (
	"fmt"
	"library/Entity"
	"library/utill"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func Login(c *gin.Context) {
	var login Entity.LoginDetails
	if err := c.BindJSON(&login); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else {

		login.Email = strings.Trim(login.Email, " ")
		login.Email = strings.ToLower(login.Email)
		login.Password = strings.Trim(login.Password, " ")

		if !(utill.IsValidEmail(login.Email)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Email Format")
			return
		}

		var user Entity.User

		result := db.First(&user, "email = ?", login.Email)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid user")
			return

		} else {

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
			if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "Incorrect Password")
				return

			} else {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"email": user.Email,
					"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
				})

				// Sign and get the complete encoded token as a string using the secret
				tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "Failed to create token",
					})
					return
				}

				// Respond
				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

				c.JSON(http.StatusOK, gin.H{"Token": tokenString, "User": user})

			}
		}

	}

}

func AuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/login") {
			c.Next()
			return
		}
		token2 := c.Request.Header.Get("Authorization")
		tokenString, err := c.Cookie("Authorization")

		if tokenString == "" || err != nil {
			if token2 != "" {
				tokenString = token2
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthenticated")
				return
			}
		}
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			email := claims["email"].(string)

			var user Entity.User
			result := db.Preload(clause.Associations).First(&user, "email = ?", email)
			if result.RowsAffected == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
				return
			} else {
				if user.Role != role {
					c.AbortWithStatusJSON(http.StatusForbidden, "Role mismatch you are not "+role+" of library ")
					c.Abort()
					return
				}
			}
			if user.Email == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Attach the request
			c.Set("user", user)

			c.Next()
		}
	}
}

func SignupReader(c *gin.Context) {
	var userDetails struct {
		Name          string
		Email         string `json:"email"`
		ContactNumber int    `json:"contact_number"`
		Password      string
	}
	lib := c.Param("LibId")
	var library Entity.Library
	var user Entity.User
	if err := c.BindJSON(&userDetails); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {

		userDetails.Email = strings.Trim(userDetails.Email, " ")
		userDetails.Password = strings.Trim(userDetails.Password, " ")
		userDetails.Name = strings.Trim(userDetails.Name, " ")
		userDetails.Email = strings.ToLower(userDetails.Email)

		fmt.Print(userDetails)
		if !(utill.IsValidEmail(userDetails.Email)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Email Format")
			return
		}
		if !(utill.IsContactNumberValid(userDetails.ContactNumber)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Contact number format enter without +91 and must be 10 number only")
			return
		}
		if !(utill.IsNameValid(userDetails.Name)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Name Format")
			return
		}
		if !(utill.IsPasswordValid(userDetails.Password)) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "password must be between 8 to 15 characters included capital letters smaller letters symbols and numbers ")
			return
		}

		result := db.Preload(clause.Associations).Find(&library, "id = ?", lib)
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, "Library not exists")
			return

		} else {
			user.Name = userDetails.Name
			user.ContactNumber = userDetails.ContactNumber
			user.Email = userDetails.Email
			user.Role = "reader"
			user.LibID = int(library.ID)
			user.Library = library
			bs, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)

			}
			user.Password = string(bs)
			tx := db.Begin()
			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"msg": "User already exists"})
				return
			} else {
				tx.Commit()
				c.JSON(http.StatusCreated, user.Name+" Added as a reader in "+library.Name)

			}
		}
	}
}
func Logout(c *gin.Context) {

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, "user logged out")

}
func Library(c *gin.Context) {
	var library []Entity.Library

	result := db.Find(&library)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, "Library Not Available at this moment")
		return
	} else {
		c.JSON(http.StatusOK, library)

	}
}
