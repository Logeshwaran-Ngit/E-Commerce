package routes

import (
	"myapp/go-gin-ecom/controllers"
	"myapp/go-gin-ecom/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	r.POST("/users", controllers.Signup)
	r.PUT("/users", middleware.RoleAuthorization("admin", "user", "seller", "distributor"), controllers.UpdatebyId)
	r.POST("/userlogin", controllers.User_signin)
	r.POST("/cart", middleware.RoleAuthorization("admin", "user"), controllers.AddToCart)
	r.GET("/cart/:id", middleware.RoleAuthorization("admin", "user"), controllers.GetUserbyID)
	r.GET("/cart/:id", middleware.RoleAuthorization("admin", "user"), controllers.GetUserCart)
	r.PUT("/cart/:id", middleware.RoleAuthorization("admin", "user"), controllers.UpdateCartStock)
	r.DELETE("/cart/:id", middleware.RoleAuthorization("admin", "user"), controllers.RemoveFromCart)

}
