package routes

import (
	"API/controller"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/", controller.GetUsers)
	router.GET("/:id", controller.GetUserID)
	router.POST("/", controller.CreateUser)
	router.DELETE("/:id", controller.DeleteUser)
	router.PUT("/:id", controller.UpdateUser)
	router.POST("/login", controller.Login) // เพิ่มบรรทัดนี้
}
