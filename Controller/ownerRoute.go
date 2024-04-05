package controller

import (
	"library/service"

	"github.com/gin-gonic/gin"
)

func ownerRoute(owner *gin.RouterGroup) {
	owner.Use(service.AuthMiddleware("owner"))
	{

		owner.GET("/user", service.GetUsers)
		owner.POST("/admin", service.AddAdmin)
		owner.DELETE("/admin", service.RemoveAdmin)
		owner.GET("/seerequest")

	}
}
