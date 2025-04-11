package handlers

import (
	"net/http"
	"server/models"
	"strconv"
)

func (h *HandlerContext) GetOrders(w http.ResponseWriter, r *http.Request) {
	userIdQuery := r.URL.Query().Get("userId")
	userId, err := strconv.ParseUint(userIdQuery, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid userId")
		return
	}

	var orders []models.Order
	if err := h.db.Preload("Items").Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch orders: "+err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Orders retrieved successfully",
		"orders":  orders,
	})
}
