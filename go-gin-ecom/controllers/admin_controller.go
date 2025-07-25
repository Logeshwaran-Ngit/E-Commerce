package controllers

import (
	"myapp/go-gin-ecom/config"
	"myapp/go-gin-ecom/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteUserbyId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if result := config.DB.Delete(&models.Users{}, id); result.RowsAffected == 0 {
		c.JSON(400, gin.H{"message": "user not avilable"})
		return
	}
	c.JSON(200, gin.H{"message": "user deleted successfuly"})
}
func GetUserbyID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	var user models.Users
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}
func GetAllUser(c *gin.Context) {
	var user []models.Users
	if err := config.DB.Find(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}
