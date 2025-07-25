package routes

import (
	"myapp/go-gin-ecom/controllers"
	"myapp/go-gin-ecom/middleware"

	"github.com/gin-gonic/gin"
)

func RegisteradminRoutes(r *gin.Engine) {
	r.DELETE("/admin/:id", middleware.RoleAuthorization("admin"), controllers.DeleteUserbyId)
	r.GET("/admin/:id", middleware.RoleAuthorization("admin"), controllers.GetUserbyID)
	r.GET("/admin", middleware.RoleAuthorization("admin"), controllers.GetAllUser)
}
