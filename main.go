package main

import (
	"API/config"
	"API/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	config.Connection()
	config.InitRedis()
	routes.UserRoute(router)
	routes.ProductRoute(router)
	router.Run(":8080")
}
