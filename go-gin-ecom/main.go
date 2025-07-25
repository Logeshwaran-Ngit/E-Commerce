package main

import (
	"myapp/go-gin-ecom/config"
	"myapp/go-gin-ecom/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.Connect()
	routes.RegisterProductRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.RegisteradminRoutes(r)
	routes.CartRoutes(r)
	r.Run("localhost:8080")
}
