package handlers

import (
	"encoding/json"
	"net/http"
	"server/models"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func (h *HandlerContext) GetCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var cartItems []models.Cart
	if err := h.db.Preload("Product").Preload("Product.Images").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch cart")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]any{
		"message":   "Cart retrieved successfully",
		"cartItems": cartItems,
	})
}

func (h *HandlerContext) AddToCart(w http.ResponseWriter, r *http.Request) {
	userIdQuery := r.URL.Query().Get("userId")
	userId, err := strconv.ParseUint(userIdQuery, 10, 32)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid userId")
		return
	}

	var reqBody struct {
		ProductId uint `json:"productId"`
		Quantity  int  `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if reqBody.ProductId == 0 {
		respondWithError(w, http.StatusBadRequest, "ProductId is reuqired")
		return
	}

	if reqBody.Quantity <= 0 {
		respondWithError(w, http.StatusBadRequest, "Quantity is required")
		return
	}

	var product models.Product
	if err := h.db.First(&product, reqBody.ProductId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error finding product")
		}
		return
	}

	if product.Stock < reqBody.Quantity {
		respondWithError(w, http.StatusBadRequest, "Insufficient stock")
		return
	}

	var cartItem models.Cart
	err = h.db.Where("user_id = ? AND product_id = ?", userId, reqBody.ProductId).First(&cartItem).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respondWithError(w, http.StatusInternalServerError, "Error checking cart")
		return
	}

	// If never add this products
	if err == gorm.ErrRecordNotFound {
		cartItem = models.Cart{
			UserID:    uint(userId),
			ProductID: reqBody.ProductId,
			Quantity:  reqBody.Quantity,
		}

		if err := h.db.Create(&cartItem).Error; err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to add to cart")
			return
		}

	} else {

		// If already exists this product then update the new quantity
		newQuantity := cartItem.Quantity + reqBody.Quantity
		if product.Stock < newQuantity {
			respondWithError(w, http.StatusBadRequest, "Insufficient stock")
			return
		}
		cartItem.Quantity = newQuantity
		if err := h.db.Save(&cartItem).Error; err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update cart")
			return
		}
	}

	respondWithJson(w, http.StatusOK, map[string]any{
		"message":  "Products add to cart successfully",
		"cartItem": cartItem,
	})
}

// UpdateCartItem - Update the quantity of a cart item
func (h *HandlerContext) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userIdQuery := r.URL.Query().Get("userId")
	userId, err := strconv.ParseUint(userIdQuery, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid userId")
		return
	}

	cartIdStr := chi.URLParam(r, "id")
	cartId, err := strconv.ParseUint(cartIdStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	var reqBody struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if reqBody.Quantity <= 0 {
		respondWithError(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	// Fetch the cart item
	var cartItem models.Cart
	if err := h.db.Preload("Product").Where("user_id = ? AND id = ?", userId, cartId).First(&cartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Cart item not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error fetching cart item")
		}
		return
	}

	// Check stock availability
	if cartItem.Product.Stock < reqBody.Quantity {
		respondWithError(w, http.StatusBadRequest, "Insufficient stock")
		return
	}

	// Update quantity
	cartItem.Quantity = reqBody.Quantity
	if err := h.db.Save(&cartItem).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update cart item")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message":  "Cart item updated successfully",
		"cartItem": cartItem,
	})
}

// RemoveFromCart - Remove an item from the cart
func (h *HandlerContext) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userIdQuery := r.URL.Query().Get("userId")
	userId, err := strconv.ParseUint(userIdQuery, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid userId")
		return
	}

	cartIdStr := chi.URLParam(r, "id")
	cartId, err := strconv.ParseUint(cartIdStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	// Fetch the cart item
	var cartItem models.Cart
	if err := h.db.Where("user_id = ? AND id = ?", userId, cartId).First(&cartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Cart item not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error fetching cart item")
		}
		return
	}

	// Delete the cart item
	if err := h.db.Delete(&cartItem).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to remove cart item")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]any{
		"message": "Cart item removed successfully",
	})
}
