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
	err := config.DB.AutoMigrate(
		&models.Users{},        // Create users table first
		&models.Product{},      // Then products
		&models.Add_cart{},     // Then add_carts (depends on both above)
		&models.ProductOrder{}, // Then product_orders
		&models.OrderProduct{}, // Then order_products
	)

	if err != nil {
		log.Fatal("Failed to auto-migrate:", err)
	}
	routes.RegisterProductRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.RegisteradminRoutes(r)
	routes.CartRoutes(r)
	r.Run("localhost:8080")
}
