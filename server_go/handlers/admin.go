package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/models"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// AdminMiddleware - Ensures user is an admin
func (h *HandlerContext) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("userId")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		var user models.User
		if err := h.db.First(&user, userID).Error; err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		if user.Role != models.RoleAdmin {
			respondWithError(w, http.StatusForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetDashboardData - Fetches dashboard metrics and data
func (h *HandlerContext) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	var totalUsers int64
	var totalOrders int64
	var totalRevenue int64
	var lowStockProducts int64

	// Fetch total users
	if err := h.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		log.Printf("Error fetching user count: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch user count")
		return
	}

	// Fetch total orders
	if err := h.db.Model(&models.Order{}).Count(&totalOrders).Error; err != nil {
		log.Printf("Error fetching order count: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch order count")
		return
	}

	// Fetch total revenue (sum of order totals)
	if err := h.db.Model(&models.Order{}).Select("COALESCE(SUM(total), 0)").Scan(&totalRevenue).Error; err != nil {
		log.Printf("Error fetching revenue: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch revenue")
		return
	}

	// Fetch low stock products (stock <= 5)
	if err := h.db.Model(&models.Product{}).Where("stock <= ?", 5).Count(&lowStockProducts).Error; err != nil {
		log.Printf("Error fetching low stock products: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch low stock products")
		return
	}

	// Fetch recent orders (last 5)
	var recentOrders []models.Order
	if err := h.db.Preload("User").Order("created_at DESC").Limit(5).Find(&recentOrders).Error; err != nil {
		log.Printf("Error fetching recent orders: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch recent orders")
		return
	}

	// Fetch all products
	var products []models.Product
	if err := h.db.Preload("Images").Preload("Category").Find(&products).Error; err != nil {
		log.Printf("Error fetching products: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	// Fetch all users
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"metrics": map[string]interface{}{
			"totalUsers":       totalUsers,
			"totalOrders":      totalOrders,
			"totalRevenue":     totalRevenue,
			"lowStockProducts": lowStockProducts,
		},
		"recentOrders": recentOrders,
		"products":     products,
		"users":        users,
	})
}

// UpdateUserRole - Updates a user's role
func (h *HandlerContext) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var reqBody struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if reqBody.Role != "user" && reqBody.Role != "admin" {
		respondWithError(w, http.StatusBadRequest, "Invalid role")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "User not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error fetching user")
		}
		return
	}

	user.Role = models.UserRole(reqBody.Role)
	if err := h.db.Save(&user).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "User role updated successfully",
	})
}

// UpdateProduct - Updates a product
func (h *HandlerContext) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid product ID: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var reqBody struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		CategoryID  uint    `json:"categoryId"`
		Stock       int     `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var product models.Product
	if err := h.db.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			log.Printf("Error fetching product: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Error fetching product")
		}
		return
	}

	product.Name = reqBody.Name
	product.Price = reqBody.Price
	product.Description = reqBody.Description
	product.CategoryID = reqBody.CategoryID
	product.Stock = reqBody.Stock

	if err := h.db.Save(&product).Error; err != nil {
		log.Printf("Error updating product: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Product updated successfully",
	})
}

// DeleteProduct - Deletes a product
func (h *HandlerContext) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid product ID: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := h.db.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			log.Printf("Error fetching product: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Error fetching product")
		}
		return
	}

	if err := h.db.Delete(&product).Error; err != nil {
		log.Printf("Error deleting product: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Product deleted successfully",
	})
}

// UpdateOrderStatus - Updates an order's status
func (h *HandlerContext) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid order ID: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var reqBody struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	validStatuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}
	isValidStatus := false
	for _, s := range validStatuses {
		if reqBody.Status == s {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		log.Printf("Invalid status: %s", reqBody.Status)
		respondWithError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	var order models.Order
	if err := h.db.First(&order, orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			log.Printf("Error fetching order: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Error fetching order")
		}
		return
	}

	order.Status = reqBody.Status
	if err := h.db.Save(&order).Error; err != nil {
		log.Printf("Error updating order status: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Order status updated successfully",
	})
}


// GetOrder - Fetches a single order with details
func (h *HandlerContext) GetOrder(w http.ResponseWriter, r *http.Request) {
    orderIDStr := chi.URLParam(r, "id")
    orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
    if err != nil {
        log.Printf("Invalid order ID: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid order ID")
        return
    }

    var order models.Order
    if err := h.db.
        Preload("Items").
        Preload("User").
        First(&order, orderID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            respondWithError(w, http.StatusNotFound, "Order not found")
        } else {
            log.Printf("Error fetching order with ID %d: %v", orderID, err)
            respondWithError(w, http.StatusInternalServerError, "Error fetching order")
        }
        return
    }

    respondWithJson(w, http.StatusOK, map[string]interface{}{
        "order": order,
    })
}

// GetAllOrders - Fetches a paginated list of orders
func (h *HandlerContext) GetAllOrders(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }

    offset := (page - 1) * limit

    var orders []models.Order
    if err := h.db.
        Preload("User").
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&orders).Error; err != nil {
        log.Printf("Error fetching orders: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Error fetching orders")
        return
    }

    respondWithJson(w, http.StatusOK, map[string]interface{}{
        "orders": orders,
    })
}

func (h *HandlerContext) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid order ID: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var reqBody struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var order models.Order
	if err := h.db.First(&order, orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			log.Printf("Error fetching order: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Error fetching order")
		}
		return
	}

	order.Status = reqBody.Status
	if err := h.db.Save(&order).Error; err != nil {
		log.Printf("Error updating order: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Order updated successfully",
	})
}
