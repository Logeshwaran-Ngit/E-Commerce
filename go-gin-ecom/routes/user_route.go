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
}
