package main

import (
	"Stant/LestaGamesInternship/internal/api/web"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"net/http"
)

func main() {
	inMemoryTermStore := stores.NewInMemoryTermStore()

	router := http.NewServeMux()
	router.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	router.Handle("GET /", web.HandleIndexGet(inMemoryTermStore))
	router.Handle("POST /", web.HandleIndexPost(inMemoryTermStore))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
