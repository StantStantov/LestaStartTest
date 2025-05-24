package main

import (
	"Stant/LestaGamesInternship/internal"
	"Stant/LestaGamesInternship/internal/stores"
	"net/http"
)

func main() {
	inMemoryTermStore := stores.NewInMemoryTermStore()

	router := http.NewServeMux()
	router.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	router.Handle("GET /", internal.HandleIndexGet(inMemoryTermStore))
	router.Handle("POST /", internal.HandleIndexPost(inMemoryTermStore))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
