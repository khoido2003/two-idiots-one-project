package main

import (
	"log"
	"net/http"
	"os"
	"server/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()
	portStr := string(os.Getenv("PORT"))

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment variable!")
		return
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{}, &models.Category{}, &models.Product{},
		&models.ProductImage{}, &models.ProductSpec{}, &models.Cart{},
		&models.Order{}, &models.OrderItem{}, &models.Review{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}

	log.Printf("Server starting on port: %v", portStr)
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal("Something went wrong: ", err)
	}
}
