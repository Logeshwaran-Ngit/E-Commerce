package routes

import (
	"myapp/go-gin-ecom/controllers"
	"myapp/go-gin-ecom/middleware"

	"github.com/gin-gonic/gin"
)

func CartRoutes(r *gin.Engine) {
	r.POST("/cart", middleware.RoleAuthorization("admin", "user"), controllers.AddToCart)
	r.GET("/cart/:id", middleware.RoleAuthorization("admin", "user"), controllers.GetUserbyID)
	r.GET("/cart", middleware.RoleAuthorization("admin", "user"), controllers.GetUserCart)
	r.PUT("/cart/:id", middleware.RoleAuthorization("admin", "user"), controllers.UpdateCartStock)
	r.DELETE("/cart/:id", middleware.RoleAuthorization("admin", "user"), controllers.RemoveFromCart)
}
