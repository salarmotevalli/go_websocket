package main

import (
	"net/http"
	"ws/internal/handlers"

	"github.com/gorilla/mux"
)

func routes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/ws", handlers.WsEndpoint)
	return mux
}
