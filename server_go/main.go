package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/handlers"
	"server/models"
	"server/routers"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Run project: go run .
// Clean up unused dependency: go mod tidy
// Dowload dependency: go mod download

func main() {

	godotenv.Load()
	portStr := string(os.Getenv("PORT"))

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment variable!")
		return
	}

	// Open connect to Postgres
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
	}

	if os.Getenv("MIGRATE") == "true" {
		fmt.Println("Start migrating database schema....")
		err = db.AutoMigrate(
			&models.User{}, &models.Category{}, &models.Product{},
			&models.ProductImage{}, &models.ProductSpec{}, &models.Cart{},
			&models.Order{}, &models.OrderItem{}, &models.Review{},
		)
		if err != nil {
			log.Fatalf("Failed to migrate: %v", err)
		}
	}

	// Create Handler Context
	handlerContext := handlers.New(db)

	// Setup router
	router := routers.NewRouter(handlerContext)

	// Init server
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
