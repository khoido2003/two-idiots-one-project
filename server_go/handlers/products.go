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

	// Ensure at least one image is primary
	hasPrimary := false
	for _, img := range req.Images {
		if img.IsPrimary {
			hasPrimary = true
			break
		}
	}
	if !hasPrimary {
		respondWithError(w, http.StatusBadRequest, "At least one image must be marked as primary")
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

// GetProducts - Retrieve products with filtering, sorting, and pagination
func (h *HandlerContext) GetProducts(w http.ResponseWriter, r *http.Request) {
	query := h.db.Session(&gorm.Session{PrepareStmt: false}).Preload("Images").Preload("Specs").Preload("Category")

	// Search filter
	if search := r.URL.Query().Get("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", searchTerm, searchTerm)
	}

	// Category filter
	if category := r.URL.Query().Get("category"); category != "" && category != "All" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.name = ?", category)
	}

	// Price range filter
	if minPriceStr := r.URL.Query().Get("minPrice"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			query = query.Where("products.price >= ?", minPrice)
		}
	}
	if maxPriceStr := r.URL.Query().Get("maxPrice"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			query = query.Where("products.price <= ?", maxPrice)
		}
	}

	// Stock filter (for limited deals)
	if maxStockStr := r.URL.Query().Get("maxStock"); maxStockStr != "" {
		if maxStock, err := strconv.Atoi(maxStockStr); err == nil {
			query = query.Where("products.stock <= ?", maxStock)
		}
	}

	// Sorting (e.g., for new arrivals)
	if sort := r.URL.Query().Get("sort"); sort != "" {
		switch sort {
		case "createdAtDesc":
			query = query.Order("products.created_at DESC")

		case "createdAtAsc":
			query = query.Order("products.created_at ASC")

		case "priceAsc":
			query = query.Order("products.price ASC")
		case "priceDesc":
			query = query.Order("products.price DESC")
		}
	}

	// Pagination
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 8
	}

	offset := (page - 1) * limit

	// Get total count for pagination
	var total int64
	countQuery := h.db.Model(&models.Product{}).Session(&gorm.Session{PrepareStmt: false})
	// Apply same filters to count query (search, category, price, stock)
	if search := r.URL.Query().Get("search"); search != "" {
		searchTerm := "%" + search + "%"
		countQuery = countQuery.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}
	if category := r.URL.Query().Get("category"); category != "" && category != "All" {
		countQuery = countQuery.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.name = ?", category)
	}
	if minPriceStr := r.URL.Query().Get("minPrice"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			countQuery = countQuery.Where("price >= ?", minPrice)
		}
	}
	if maxPriceStr := r.URL.Query().Get("maxPrice"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			countQuery = countQuery.Where("price <= ?", maxPrice)
		}
	}
	if maxStockStr := r.URL.Query().Get("maxStock"); maxStockStr != "" {
		if maxStock, err := strconv.Atoi(maxStockStr); err == nil {
			countQuery = countQuery.Where("stock <= ?", maxStock)
		}
	}
	if err := countQuery.Count(&total).Error; err != nil {
		log.Printf("Failed to count products: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to count products")
		return
	}

	// Fetch paginated products
	var products []models.Product
	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		log.Printf("Failed to fetch products: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message":     "Products retrieved successfully",
		"products":    products,
		"totalPages":  totalPages,
		"currentPage": page,
		"totalItems":  total,
	})
}

// GetProductByID - Retrieve a specific product by ID with related products
func (h *HandlerContext) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Fetch the product with preloaded data
	var product models.Product
	if err := h.db.Preload("Images").Preload("Specs").Preload("Category").First(&product, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			log.Printf("Failed to fetch product %d: %v", id, err)
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch product")
		}
		return
	}

	// Fetch related products (e.g., same category, exclude current product, limit 3)
	var related []models.Product
	if err := h.db.
		Preload("Images").
		Where("category_id = ? AND id != ?", product.CategoryID, product.ID).
		Limit(3).
		Find(&related).Error; err != nil {
		log.Printf("Failed to fetch related products for product %d: %v", id, err)
		// Don't fail the request, just log and proceed without related products
	}

	// Populate RelatedIDs
	relatedIDs := make([]uint, len(related))
	for i, rel := range related {
		relatedIDs[i] = rel.ID
	}

	// Response structure with related IDs
	type ProductResponse struct {
		ID          uint                  `json:"id"`
		Name        string                `json:"name"`
		Price       float64               `json:"price"`
		Description string                `json:"description"`
		Category    models.Category       `json:"category"`
		Stock       int                   `json:"stock"`
		Images      []models.ProductImage `json:"images"`
		Specs       []models.ProductSpec  `json:"specs"`
		RelatedIDs  []uint                `json:"related"`
	}

	response := ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Category:    product.Category,
		Stock:       product.Stock,
		Images:      product.Images,
		Specs:       product.Specs,
		RelatedIDs:  relatedIDs,
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Product retrieved successfully",
		"product": response,
	})
}
