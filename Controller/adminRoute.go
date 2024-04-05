package controller

import (
	"library/service"

	"github.com/gin-gonic/gin"
)

func adminRoute(admin *gin.RouterGroup) {
	admin.Use(service.AuthMiddleware("admin"))

	{

		admin.GET("/user", service.GetUsers)
		admin.POST("/book", service.AddBook)
		admin.DELETE("/book/:ISBN", service.RemoveBook)
		admin.GET("/book/:ISBN", service.BookById)
		admin.GET("/book", service.GetAllBooks)
		admin.PUT("/book/:ISBN", service.UpdateBook)
		admin.GET("/request", service.IssueRequest)
		admin.PUT("/request/approve/:rid", service.ApproveRequest)
		admin.PUT("/request/reject/:rid", service.RejectRequest)
		admin.GET("/issued", service.IssuedBook)

		admin.GET("/issue/:id", service.UserIssue)

		admin.GET("/decreaseBook/:isbn", service.DecreaseBook)
		admin.GET("/increaseBook/:isbn", service.IncreaseBook)
		admin.POST("/invite", service.Invite)

	}
}
