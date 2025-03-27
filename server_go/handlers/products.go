package handlers

import (
	"encoding/json"
	"net/http"
	"server/models"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type CreateProductRequest struct {
	Name        string                `json:"name"`
	Price       float64               `json:"price"`
	Description string                `json:"description"`
	CategoryID  uint                  `json:"categoryId"`
	Stock       int                   `json:"stock"`
	Images      []models.ProductImage `json:"images"`
	Specs       []models.ProductSpec  `json:"specs"`
}

func (h *HandlerContext) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validation
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Product name is required")
		return
	}
	if req.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "Price must be greater than 0")
		return
	}
	if req.Description == "" {
		respondWithError(w, http.StatusBadRequest, "Description is required")
		return
	}
	if req.CategoryID == 0 {
		respondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}
	if len(req.Images) == 0 {
		respondWithError(w, http.StatusBadRequest, "At least one image is required")
		return
	}

	// Check if category exists
	var category models.Category
	if err := h.db.First(&category, req.CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusBadRequest, "Category not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error checking category")
		}
		return
	}

	// Create product
	product := models.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Stock:       req.Stock,
		Images:      req.Images,
		Specs:       req.Specs,
	}

	if err := h.db.Create(&product).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	respondWithJson(w, http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"product": product,
	})
}

// GetProducts - Retrieve all products
func (h *HandlerContext) GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	// Preload related data (Images, Specs, Category)
	if err := h.db.Preload("Images").Preload("Specs").Preload("Category").Find(&products).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message":  "Products retrieved successfully",
		"products": products,
	})
}

// GetProductByID - Retrieve a specific product by ID
func (h *HandlerContext) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	// Preload related data
	if err := h.db.Preload("Images").Preload("Specs").Preload("Category").First(&product, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch product")
		}
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Product retrieved successfully",
		"product": product,
	})
}
