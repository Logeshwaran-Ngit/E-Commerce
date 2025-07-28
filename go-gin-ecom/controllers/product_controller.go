package controllers

import (
	"fmt"
	"math/rand"
	"myapp/go-gin-ecom/config"
	"myapp/go-gin-ecom/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Add_product(c *gin.Context) {
	var product models.Product

	// Bind the incoming JSON to the product struct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Check if the same product with name & model already exists
	var existingProduct models.Product
	err := config.DB.Where("name = ? AND model = ?", product.Name, product.Model).First(&existingProduct).Error
	if err == nil {
		// Product exists, just update the stock
		existingProduct.Stock += product.Stock
		if err := config.DB.Save(&existingProduct).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":         "Stock updated successfully",
			"updated_product": existingProduct,
		})
		return
	}

	// Product doesn't exist, create a new one
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Product created successfully",
		"created_product": product,
	})
}

func Get_product(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product fetched successfully", "product": product})
}

func Remove_product(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.Delete(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func View_all_product(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

type ProductItem struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type ProductOrderRequest struct {
	UserID   uint          `json:"user_id"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	PhoneNo  string        `json:"phone_no"`
	Address  string        `json:"address"`
	Products []ProductItem `json:"products"`
}

func PlaceOrder(c *gin.Context) {
	var req ProductOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID == 0 || len(req.Products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and products are required"})
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// 🔄 Remove old incomplete "Not Placed" orders for this user
	if err := tx.Where("user_id = ? AND status = ?", req.UserID, "Not Placed").Delete(&models.ProductOrder{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clean old orders"})
		return
	}

	var orderProducts []models.OrderProduct
	totalAmount := 0.0

	for _, item := range req.Products {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product ID %d not found", item.ProductID)})
			return
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Not enough stock for %s", product.Name)})
			return
		}

		subTotal := float64(product.Price) * float64(item.Quantity)
		totalAmount += subTotal

		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
			return
		}

		orderProducts = append(orderProducts, models.OrderProduct{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  product.Price,
			TotalPrice: subTotal,
		})
	}

	orderID := fmt.Sprintf("ORD-%06d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000))
	deliveryCharge := 50.0

	order := models.ProductOrder{
		OrderID:        orderID,
		UserID:         req.UserID,
		Name:           req.Name,
		Email:          req.Email,
		PhoneNo:        req.PhoneNo,
		Address:        req.Address,
		CurrentAmount:  totalAmount,
		DeliveryCharge: deliveryCharge,
		TotalAmount:    totalAmount + deliveryCharge,
		Status:         "Placed",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order"})
		return
	}

	// Set ProductOrderID for all items
	for i := range orderProducts {
		orderProducts[i].ProductOrderID = order.ID
	}

	if err := tx.Create(&orderProducts).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order items"})
		return
	}

	tx.Commit()

	order.OrderItems = orderProducts

	c.JSON(http.StatusOK, gin.H{
		"message": "Order placed successfully",
		"order":   order,
	})
}

func Pending_Delivery(c *gin.Context) {
	var orders []models.ProductOrder

	if err := config.DB.
		Where("status = ?", "Not Placed").
		Preload("OrderItems").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending deliveries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fetched all pending deliveries",
		"orders":  orders,
	})
}

func Pending_Delivery_user(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var orders []models.ProductOrder
	if err := config.DB.Preload("OrderItems").Where("user_id = ? AND status = ?", userID, "Not Placed").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's pending deliveries"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No pending deliveries found for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fetched user's pending deliveries", "orders": orders})
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
	if product.Stock < input.Product_Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	var existingCart models.Add_cart
	err := config.DB.Where("user_id = ? AND product_id = ?", input.User_Id, input.Product_Id).First(&existingCart).Error
	if err == nil {
		existingCart.Product_Stock += input.Product_Stock
		config.DB.Save(&existingCart)
		c.JSON(http.StatusOK, gin.H{"message": "Cart updated with additional quantity"})
		return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var cart models.Add_cart
	if err := config.DB.First(&cart, cartID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	cart.Product_Stock = input.NewStock
	if err := config.DB.Save(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart stock updated successfully"})
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
	if err := config.DB.Preload("Product").Preload("User").Where("user_id = ?", userID).Find(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart for user"})
		return
	}

	if len(cart) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No cart items found for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User cart fetched successfully", "cart": cart})
}
