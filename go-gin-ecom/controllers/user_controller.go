package controllers

import (
	"myapp/go-gin-ecom/config"
	"myapp/go-gin-ecom/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existingUser models.Users
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}
	rotationStr := os.Getenv("ROTATION")
	rotation, err := strconv.Atoi(rotationStr)
	if err != nil {
		rotation = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), rotation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hash)
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}
func User_signin(c *gin.Context) {
	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	var login Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.Users
	if err := config.DB.First(&user, "email=?", login.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalied email id and possword"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		c.JSON(409, gin.H{"error": "invalied password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"Role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 1).Unix()})
	tokenstring, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(403, gin.H{"error": "could not create a token"})
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenstring, 3600*24*1, "/", "", false, true)
	c.JSON(200, gin.H{"message": "login successfuly"})
}
func UpdatebyId(c *gin.Context) {
	var update models.Users
	if err := c.ShouldBindJSON(&update); err != nil || update.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input or missing user ID"})
		return
	}
	var user models.Users
	if err := config.DB.First(&user, update.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Name = update.Name
	user.Email = update.Email
	user.Role = update.Role
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}
func Add_to_cart(c *gin.Context) {
	var new_cart models.Add_cart
	if err := c.ShouldBindJSON(&new_cart); err != nil {
		c.JSON(400, gin.H{"error": "json file not valied"})
	}
	err := config.DB.Where("ID = ? AND Product_id = ?", new_cart.User_Id, new_cart.Product_Id).First(&new_cart).Error
	if err != nil {
		c.JSON(400, gin.H{"error": "it no product or id is not avliable"})
	}
}
