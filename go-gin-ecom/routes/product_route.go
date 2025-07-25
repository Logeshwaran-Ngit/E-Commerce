package routes

import (
	"myapp/go-gin-ecom/controllers"
	"myapp/go/gin-ecom/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine) {
	r.POST("/products", middleware.RoleAuthorization("admin", "seller"), controllers.Add_product)
	r.GET("/products", middleware.RoleAuthorization("admin", "user", "seller"), controllers.View_all_product)
	r.GET("/products/:id", middleware.RoleAuthorization("admin", "user", "seller"), controllers.Get_product)
	r.DELETE("/products/:id", middleware.RoleAuthorization("admin", "seller"), controllers.Remove_product)
	r.POST("/place-order", middleware.RoleAuthorization("admin", "user"), controllers.PlaceOrder)
	r.GET("/pending-product", middleware.RoleAuthorization("admin", "Delivery_Associate", "seller"), controllers.Pending_Delivery)
	r.GET("/pending", middleware.RoleAuthorization("admin", "Delivery_Associate"), controllers.Pending_Delivery)
	r.GET("/pending/:id", middleware.RoleAuthorization("admin", "Delivery_Associate"), controllers.Pending_Delivery_user)
}
