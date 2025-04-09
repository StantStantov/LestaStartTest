package main

import (
	"Stant/LestaGamesInternship/internal"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	router.Handle("GET /", internal.HandleIndexGet())
	router.Handle("POST /", internal.HandleIndexPost())

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
