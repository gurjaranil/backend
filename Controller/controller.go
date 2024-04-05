package controller

import (
	"library/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "token", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/login", service.Login)
	r.POST("/Register/Library", service.RegisterLibrary)
	r.POST("/signup/:LibId", service.SignupReader)
	r.GET("/logout", service.Logout)
	r.GET("/library", service.Library)

	owner := r.Group("/owner")
	ownerRoute(owner)

	admin := r.Group("/admin")
	adminRoute(admin)

	reader := r.Group("/reader")
	readerRoute(reader)

	return r
}
