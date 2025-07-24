package routes

import (
	"github.com/gin-gonic/gin"
)

func AdminRoute(r *gin.Engine) {
	r.POST("/user")
}
