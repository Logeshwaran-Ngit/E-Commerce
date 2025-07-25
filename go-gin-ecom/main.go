package main

import (
	"log"
	"myapp/go-gin-ecom/config"
	"myapp/go-gin-ecom/models"
	"myapp/go-gin-ecom/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.Connect()
	err := config.DB.AutoMigrate(&models.Users{}, &models.Product{}, &models.Add_cart{}, &models.ProductOrder{}, &models.OrderProduct{})
	if err != nil {
		log.Fatal("Failed to auto-migrate:", err)
	}
	routes.RegisterProductRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.RegisteradminRoutes(r)
	routes.CartRoutes(r)
	r.Run("localhost:8080")
}
