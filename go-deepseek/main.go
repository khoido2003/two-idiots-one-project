package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/ollama/ollama/api"
)

// Test server: npx wscat -c ws://127.0.0.1:9089/ws
const PORT = ":9089"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketHandler(w http.ResponseWriter, r *http.Request, client *api.Client) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error starting socket:  ", err)
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
			log.Println(mes)

			go processMessage(conn, client, &mes)
		}
	}()

}

func processMessage(conn *websocket.Conn, client *api.Client, mes *string) {

	// Prepare ollama request form
	ollamaMessages := []api.Message{
		{
			Role:    "system",
			Content: "You are a helpful store assistant. Provide concise response",
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

	respFunc := func(fullResponse string) error {
		// Sending back the response from the model to the websocket
		log.Printf("%v", fullResponse)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(fullResponse)); err != nil {
			log.Println("Error writing message: ", err)
			return err
		}
		return nil
	}
	// Buffer to collect the full response
	var fullResponse string

	// Send request to the model and collect the full response
	err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		fullResponse += resp.Message.Content
		return nil
	})

	err = respFunc(fullResponse)
	if err != nil {
		log.Println("Ollama chat error: ", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error generating response by Deepseek ollama"))
	}
}

func main() {

	// Channel to handle error
	serverError := make(chan error, 1)

	// Create Ollama client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Println("Ollama client error:", err)
		serverError <- err
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socketHandler(w, r, client)
	})

	// Start the server
	go func() {
		log.Println("Server is running on: ", PORT)

		// Blocking loop function
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
		log.Println("Receive signal to shutdown server: ", shutdown)
	}
}
