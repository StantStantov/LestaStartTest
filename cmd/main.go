package main

import (
	"Stant/LestaGamesInternship/internal/api/rest"
	"Stant/LestaGamesInternship/internal/api/web"
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"fmt"
	"log"
	"net/http"
)

func main() {
	appConfig, err := config.ReadAppConfig()
	if err != nil {
		log.Printf("cmd/main.main: [%v]", err)
		return
	}

	inMemoryTermStore := stores.NewInMemoryTermStore()
	inMemoryMetricStore := stores.NewInMemoryMetricStore()

	router := http.NewServeMux()
	web.SetupWebRouter(router, inMemoryMetricStore, inMemoryTermStore)
	rest.SetupRestRouter(router, inMemoryMetricStore)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.ServerPort()),
		Handler: router,
	}

	server.ListenAndServe()
}
