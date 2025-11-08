package main

import (
	"API/config"
	"API/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	config.Connection()
	routes.UserRoute(router)
	router.Run(":8080")
}
