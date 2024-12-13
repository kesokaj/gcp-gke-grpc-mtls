package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	ENV_GRPC_PORT = getEnv("GRPC_PORT", "50051")
	ENV_WEB_PORT  = getEnv("WEB_PORT", "8080")
)

func main() {
	// Get the port to check from environment variable, defaulting to 50051

	// Create a health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Attempt to connect to the gRPC server
		conn, err := net.Dial("tcp", "localhost:"+ENV_GRPC_PORT)
		if err != nil {
			log.Printf("Error connecting to gRPC server: %v", err)
			http.Error(w, "gRPC server not available", http.StatusServiceUnavailable)
			return
		}
		defer conn.Close()

		// Server is reachable, respond with 200 OK and JSON message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Start the HTTP server
	log.Printf("Health check server listening on port " + ENV_WEB_PORT)
	log.Fatal(http.ListenAndServe(":"+ENV_WEB_PORT, nil))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
