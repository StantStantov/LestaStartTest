package main

import (
	"Stant/LestaGamesInternship/internal/api/rest"
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	appConfig, err := config.ReadAppConfig()
	if err != nil {
		log.Printf("cmd/main.main: [%v]", err)
		return
	}

	dbPool, err := pgxpool.New(context.Background(), appConfig.DatabaseUrl())
	if err != nil {
		log.Printf("cmd/main.main: [%v]", err)
		return
	}

	metricStore := pgsql.NewMetricStore(dbPool)

	router := http.NewServeMux()
	rest.SetupRestRouter(router, metricStore)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.ServerPort()),
		Handler: router,
	}

	server.ListenAndServe()
}
