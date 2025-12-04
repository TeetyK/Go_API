package main

import (
	"API/config"
	"API/controller"
	"API/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// router := gin.New()
	config.Connection()
	config.InitRedis()
	// routes.UserRoute(router)
	// routes.ProductRoute(router)

	router := gin.Default()
	router.POST("/login", middleware.RateLimiter(), controller.Login)
	router.POST("/users", controller.CreateUser)

	authorized := router.Group("/")
	authorized.Use(middleware.RequireAuth)
	{
		authorized.GET("/users", controller.GetUsers)
		authorized.GET("/users/:id", controller.GetUserID)
		authorized.PUT("/users/:id", controller.UpdateUser)
		authorized.DELETE("/users/:id", controller.DeleteUser)
		authorized.GET("/products", controller.GetProducts)
		authorized.GET("/products/:id", controller.GetProductByID)
		authorized.POST("/products", controller.CreateProduct)
		authorized.PUT("/products/:id", controller.UpdateProduct)
		authorized.DELETE("/products/:id", controller.DeleteProduct)
	}

	router.Run(":8080")
}
