package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type HandlerContext struct {
	db gorm.DB
}

// Create new context
func New(db *gorm.DB) HandlerContext {
	return HandlerContext{
		db: *db,
	}
}

/////////////////////////////////////////////////

// Utilities method

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to encode JSON response %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithString(w http.ResponseWriter, code int, payload string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Respond with 5XX Error: %v", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}
