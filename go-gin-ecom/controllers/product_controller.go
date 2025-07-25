package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"myapp/go-gin-ecom/background_job"
	"myapp/go-gin-ecom/config"
	"myapp/go-gin-ecom/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Add_product(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var e_product models.Product
	err := config.DB.Where("Product_Name=?", "Product_Model=?", product.Product_Name, product.Product_Model).First(&e_product).Error
	if err == nil {
		e_product.Product_Stock += product.Product_Stock
		if err := config.DB.Save(&e_product); err != nil {
			c.JSON(400, gin.H{"error": "filed to update the stock"})
			return
		}
		c.JSON(200, gin.H{"message": "Stock update successfuly", "update_product": e_product})
	}
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "product created successfully", "create_product": product})
}
func Get_product(c *gin.Context) {
	num := c.Param("id")
	var product models.Product
	if err := config.DB.Find(&product, num).Error; err != nil {
		c.JSON(400, gin.H{"error": "user not avilable"})
		return
	}
	c.JSON(200, gin.H{"message": "successfuly updated ", "user": product})
}
func Remove_product(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.Delete(&product, id); err != nil {
		c.JSON(201, gin.H{"error": "the product is not avilable"})
		return
	}
	c.JSON(200, gin.H{"message": "product deleted sucessfuly"})
}
func View_all_product(c *gin.Context) {
	var product []models.Product
	if err := config.DB.Find(&product).Error; err != nil {
		c.JSON(201, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"product": product})
}

type ProductOrderRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhoneNo  string `json:"phoneno"`
	Address  string `json:"address"`
	Products []struct {
		ProductID uint `json:"product_id"`
		Quantity  uint `json:"quantity"`
	} `json:"products"`
}

func PlaceOrder(c *gin.Context) {

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var req ProductOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderProducts []models.OrderProduct
	var totalAmount float64 = 0

	for _, item := range req.Products {
		var product models.Product
		if err := config.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product ID %d not found", item.ProductID)})
			return
		}
		if product.Product_Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Not enough stock for %s", product.Product_Name)})
			return
		}

		unitPrice := float64(product.Product_Prize)
		subTotal := unitPrice * float64(item.Quantity)
		totalAmount += subTotal

		product.Product_Stock -= item.Quantity
		config.DB.Save(&product)

		orderProducts = append(orderProducts, models.OrderProduct{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  unitPrice,
			TotalPrice: subTotal,
		})
	}

	orderID := fmt.Sprintf("ORD-%d", rand.Intn(1000000))
	deliveryCharge := 50.0

	order := models.ProductOrder{
		UserID:         userID,
		OrderID:        orderID,
		Name:           req.Name,
		Email:          req.Email,
		PhoneNo:        req.PhoneNo,
		Address:        req.Address,
		CurrentAmount:  totalAmount,
		DeliveryCharge: deliveryCharge,
		TotalAmount:    totalAmount + deliveryCharge,
		Status:         "Not Placed",
		OrderItems:     []models.OrderProduct{},
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order"})
		return
	}

	for i := range orderProducts {
		orderProducts[i].OrderID = orderID
	}

	if err := config.DB.Create(&orderProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order items"})
		return
	}

	emailData := map[string]string{
		"Name":        req.Name,
		"OrderID":     orderID,
		"TotalAmount": fmt.Sprintf("%.2f", totalAmount+deliveryCharge),
		"Address":     req.Address,
	}

	go func() {
		if err := background_job.SendOrderConfirmation(req.Email, emailData); err != nil {
			log.Println("Background email error:", err)
		}
	}()
	order.OrderItems = orderProducts

	c.JSON(http.StatusOK, gin.H{
		"message": "Order placed successfully",
		"order":   order,
	})

}

func Pending_Delivery(c *gin.Context) {
	var product []models.OrderProduct
	if err := config.DB.Where("Product_Stock=?", "Not Placed").First(&product).Error; err != nil {
		c.JSON(400, gin.H{"error": "Not placed product not avilable"})
		return
	}
	c.JSON(200, gin.H{"message": "successfuly all product will be get", "product": product})
}

func Pending_Delivery_user(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var orders []models.ProductOrder
	if err := config.DB.Preload("OrderItems").
		Where("user_id = ? AND status = ?", userID, "Not Placed").
		Find(&orders).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch pending deliveries", "details": err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(404, gin.H{"error": "No pending deliveries found for this user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully fetched user's pending deliveries",
		"orders":  orders,
	})
}
func AddToCart(c *gin.Context) {
	var input models.Add_cart
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, input.Product_Id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product ID %d not found", input.Product_Id)})
		return
	}

	if product.Product_Stock < input.Product_Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	// Check if the product is already in the user's cart
	var existingCart models.Add_cart
	err := config.DB.Where("user_id = ? AND product_id = ?", input.User_Id, input.Product_Id).First(&existingCart).Error
	if err == nil {
		// Update existing quantity
		existingCart.Product_Stock += input.Product_Stock
		config.DB.Save(&existingCart)
		c.JSON(http.StatusOK, gin.H{"message": "Cart updated with more quantity"})
		return
	}

	// Else, add new cart item
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}

func UpdateCartStock(c *gin.Context) {
	cartID := c.Param("id")
	var input struct {
		NewStock uint `json:"stock"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var cart models.Add_cart
	if err := config.DB.First(&cart, cartID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cart item not found"})
		return
	}

	cart.Product_Stock = input.NewStock
	config.DB.Save(&cart)

	c.JSON(200, gin.H{"message": "Cart stock updated successfully"})
}

func RemoveFromCart(c *gin.Context) {
	productIDParam := c.Param("product_id")
	userIDParam := c.Query("user_id")

	productID, err1 := strconv.Atoi(productIDParam)
	userID, err2 := strconv.Atoi(userIDParam)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id or user_id"})
		return
	}

	var cartItem models.Add_cart
	if err := config.DB.Where("product_id = ? AND user_id = ?", productID, userID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found for this user"})
		return
	}

	if err := config.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed successfully"})
}

func GetUserCart(c *gin.Context) {
	userID := c.Query("user_id")

	var cart []models.Add_cart
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cart).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch cart for user"})
		return
	}

	if len(cart) == 0 {
		c.JSON(404, gin.H{"message": "No cart items found for this user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User cart fetched",
		"cart":    cart,
	})
}
