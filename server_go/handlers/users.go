package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"server/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	jwtSecret = "80d5eff081df7441ba1490f9ff7b5793fceaccf3e62c2159f918b6d3a2d5d2cac25f45c2382be09e83fbca4be50f56d3c004788c0bce82043efd51891fb43e0f"
)

// SignupRequest for user registration
type SignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

// SigninRequest for user login
type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup - Register a new user with default "user" role
func (h *HandlerContext) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validation
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.Phone == "" {
		respondWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		respondWithError(w, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         models.RoleUser, // Default role: user
	}

	if err := h.db.Create(&user).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT with role
	token, err := generateJWT(user.ID, string(user.Role))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJson(w, http.StatusCreated, map[string]any{
		"message": "User created successfully",
		"user": map[string]any{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"phone":     user.Phone,
		},
		"token": token,
	})
}

// Signin - Authenticate a user
func (h *HandlerContext) Signin(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT with role
	token, err := generateJWT(user.ID, string(user.Role))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Signed in successfully",
		"user": map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"role":      user.Role,
			"phone":     user.Phone,
		},
		"token": token,
	})
}

// AuthMiddleware - Validate JWT and attach user ID and role to context
func (h *HandlerContext) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}

		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		} else {
			respondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		userID, role, err := validateJWT(tokenStr)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to generate JWT with role
func generateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// Helper to validate JWT
func validateJWT(tokenStr string) (uint, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["userId"].(float64)
		role, roleOk := claims["role"].(string)
		if !ok || !roleOk {
			return 0, "", fmt.Errorf("invalid claims")
		}
		return uint(userIDFloat), role, nil
	}
	return 0, "", fmt.Errorf("invalid token")
}

// Update getUserIDFromRequest
func getUserIDFromRequest(r *http.Request) uint {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}

// ... rest of handlers (GetCart, AddToCart, etc.) unchanged, already using getUserIDFromRequest ...
