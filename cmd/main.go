package main

import (
	"Stant/LestaGamesInternship/internal"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.Handle("GET /", internal.HandleIndex())

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
