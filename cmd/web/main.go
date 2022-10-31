package main

import (
	"log"
	"net/http"
	"ws/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Start listening ws serer")
	go handlers.ListenToWsChannel()

	log.Println("server is running on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}
