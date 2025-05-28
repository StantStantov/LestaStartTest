package main

import (
	"Stant/LestaGamesInternship/internal/api/rest"
	"Stant/LestaGamesInternship/internal/api/web"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"net/http"
)

func main() {
	inMemoryTermStore := stores.NewInMemoryTermStore()
	inMemoryMetricStore := stores.NewInMemoryMetricStore()

	router := http.NewServeMux()
	web.SetupWebRouter(router, inMemoryMetricStore, inMemoryTermStore)
	rest.SetupRestRouter(router, inMemoryMetricStore)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
