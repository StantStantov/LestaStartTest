package main

import (
	"Stant/LestaGamesInternship/internal"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.Handle("GET /", internal.HandleIndexGet())
	router.Handle("POST /", internal.HandleIndexPost())

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
