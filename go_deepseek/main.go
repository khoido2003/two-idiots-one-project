package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test server: npx wscat -c ws://127.0.0.1:9089/ws
const PORT = ":9089"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Global variable to store product context
var productContext string

// InitializeDBAndFetchData connects to the database and fetches product data once
func InitializeDBAndFetchData() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return fmt.Errorf("DB_URL is not found in the environment variable")
	}

	// Open connection to Postgres
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	log.Println("Connected to DB successfully!")

	// Fetch product data with Category preloaded
	var products []Product
	result := db.Preload("Category").Find(&products)
	if result.Error != nil {
		return fmt.Errorf("failed to fetch products: %v", result.Error)
	}

	// Format product data as a string for the model
	var context string
	for _, product := range products {
		context += fmt.Sprintf(
			"Product ID: %d, Name: %s, Price: %.2f, Stock: %d, Description: %s, Category: %s\n",
			product.ID, product.Name, product.Price, product.Stock, product.Description, product.Category.Name,
		)
	}

	// Store in global variable
	productContext = context
	log.Println("Product context loaded successfully!")
	return nil
}

func socketHandler(w http.ResponseWriter, r *http.Request, client *api.Client) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error starting socket: ", err)
		return
	}

	go func() {
		defer conn.Close()
		// Receive message from client
		for {
			_, buffer, err := conn.ReadMessage()
			if err != nil {
				return
			}

			mes := string(buffer[:])
			log.Println("User message:", mes)

			go processMessage(conn, client, &mes)
		}
	}()
}

func processMessage(conn *websocket.Conn, client *api.Client, mes *string) {
	// Use pre-fetched product context (no DB query here)
	if productContext == "" {
		log.Println("Product context not initialized")
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Store data not available"))
		return
	}

	// Prepare system message with pre-fetched product context
	systemContent := fmt.Sprintf(
		"You are a helpful store assistant. Use the following product information to assist the user:\n\n%s\nProvide concise responses based on this data.",
		productContext,
	)

	// Prepare ollama request
	ollamaMessages := []api.Message{
		{
			Role:    "system",
			Content: systemContent,
		},
		{
			Role:    "user",
			Content: *mes,
		},
	}

	req := &api.ChatRequest{
		Model:    "deepseek-r1:1.5b",
		Messages: ollamaMessages,
	}

	ctx := context.Background()

	// Function to send full data back to the client
	consoleFunc := func(fullResponse string) error {
		log.Printf("Full response: %v", fullResponse)
		return nil
	}

	responseFunc := func(word string) error {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(word)); err != nil {
			log.Println("Error writing message: ", err)
			return err
		}
		return nil
	}

	// Buffer to collect the full response
	var fullResponse string

	// Send request to the model and collect the full response
	err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		err := responseFunc(resp.Message.Content)
		if err != nil {
			return err
		}
		fullResponse += resp.Message.Content
		return nil
	})

	if err != nil {
		log.Println("Ollama chat error: ", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error generating response by Deepseek ollama"))
		return
	}

	err = consoleFunc(fullResponse)
	if err != nil {
		log.Println("Console func error: ", err)
	}
}

func main() {
	// Channel to handle error
	serverError := make(chan error, 1)

	// Initialize database and fetch product data at startup
	err := InitializeDBAndFetchData()
	if err != nil {
		log.Println("Error initializing DB and fetching data:", err)
		serverError <- err
		return
	}

	// Create Ollama client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Println("Ollama client error:", err)
		serverError <- err
		return
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socketHandler(w, r, client)
	})

	// Start the server
	go func() {
		log.Println("Server is running on: ", PORT)
		err := http.ListenAndServe(PORT, nil)
		if err != nil {
			serverError <- err
		}
	}()

	// Channel to shut down server gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case shutdown := <-stop:
		log.Println("Received signal to shutdown server: ", shutdown)
	}
}
