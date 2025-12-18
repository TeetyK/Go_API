package routes

import (
	"API/controller"
	"API/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/", controller.GetUsers)
	router.GET("/:id", controller.GetUserID)
	router.POST("/", controller.CreateUser)
	router.DELETE("/:id", controller.DeleteUser)
	router.PUT("/:id", controller.UpdateUser)
	router.POST("/login", middleware.RateLimiter(), controller.Login) // เพิ่มบรรทัดนี้
	router.POST("/register", middleware.RateLimiter(), controller.CreateUser)
}
