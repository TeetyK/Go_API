package main

import (
	"API/config"
	"API/controller"
	"API/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// router := gin.New()
	config.Connection()
	config.InitRedis()
	// routes.UserRoute(router)
	// routes.ProductRoute(router)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.POST("/login", middleware.RateLimiter(), controller.Login)
	router.POST("/register", middleware.RateLimiter(), controller.CreateUser)
	router.POST("/forgot-password", middleware.RateLimiter(), controller.ForgotPassword)
	router.POST("/reset-password", middleware.RateLimiter(), controller.ResetPassword)
	// Note: The route below also creates a user, but without a rate limit.
	// Consider removing it in favor of the /register endpoint.
	// router.POST("/users", controller.CreateUser)

	authorized := router.Group("/")
	authorized.Use(middleware.RequireAuth)
	{
		authorized.GET("/", func(c *gin.Context) {
			user, exist := c.Get("user")
			if !exist {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to authorized area.",
				"user":    user,
			})
		})
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
