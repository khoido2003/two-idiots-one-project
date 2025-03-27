package handlers

import (
	"net/http"
	"server/models"
)

func (hc *HandlerContext) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	if err := hc.db.Find(&categories).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message":    "Categories retrieved successfully",
		"categories": categories,
	})

}
