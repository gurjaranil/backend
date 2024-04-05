package controller

import (
	"library/service"

	"github.com/gin-gonic/gin"
)

func readerRoute(reader *gin.RouterGroup) {
	reader.Use(service.AuthMiddleware("reader"))

	{
		reader.GET("/book/:keyword", service.SearchBook)
		reader.POST("/request", service.CreateIssue)
		reader.POST("/return/:isbn", service.ReturnBook)
		reader.GET("/book", service.GetAllBooks)
		reader.GET("/issue/:id", service.UserIssue)

	}
}
