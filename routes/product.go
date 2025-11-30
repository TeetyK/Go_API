package routes

import (
	"API/controller"
	"github.com/gin-gonic/gin"
)

// ProductRoute sets up the routes for the product resource.
func ProductRoute(router *gin.Engine) {
	// Group routes for better organization
	productRoutes := router.Group("/products")
	{
		productRoutes.GET("/", controller.GetProducts)
		productRoutes.GET("/:id", controller.GetProductByID)
		productRoutes.POST("/", controller.CreateProduct)
		productRoutes.PUT("/:id", controller.UpdateProduct)
		productRoutes.DELETE("/:id", controller.DeleteProduct)
	}
}
