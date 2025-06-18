package main

import (
	"Stant/LestaGamesInternship/internal/api/rest"
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/app/services/cllcsserv"
	"Stant/LestaGamesInternship/internal/app/services/docserv"
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"Stant/LestaGamesInternship/internal/app/services/usrserv"
	"Stant/LestaGamesInternship/internal/domain/services"
	"Stant/LestaGamesInternship/internal/infra/bcrypt"
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"Stant/LestaGamesInternship/internal/infra/snowflake"
	"Stant/LestaGamesInternship/internal/infra/volume"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	appConfig, err := config.ReadAppConfig()
	if err != nil {
		log.Fatalf("cmd/main.main: [%v]", err)
	}

	dbPool, err := pgxpool.New(context.Background(), appConfig.DatabaseUrl())
	if err != nil {
		log.Fatalf("cmd/main.main: [%v]", err)
	}
	defer dbPool.Close()

	idGenerator, err := snowflake.NewIdGenerator(1)
	if err != nil {
		log.Fatalf("cmd/main.main: [%v]", err)
	}
	passwordEncrypter := services.PasswordEncrypterFunc(bcrypt.Encrypt)
	passwordValidator := services.PasswordValidatorFunc(bcrypt.ComparePasswords)

	metricStore := pgsql.NewMetricStore(dbPool)
	sessionStore := pgsql.NewSessionStore(dbPool)
	defer sessionStore.StopCleanup(sessionStore.StartCleanup(5 * time.Minute))
	userStore := pgsql.NewUserStore(dbPool, passwordEncrypter)
	fileStore := volume.NewFileStore(appConfig.PathToDocuments())
	documentStore := pgsql.NewDocumentStore(dbPool, fileStore)
	collectionStore := pgsql.NewCollectionStore(dbPool, documentStore)

	sessionService := sesserv.NewSessionService(sessionStore,
		appConfig.ServerDomain(),
		"SESSION",
		10*time.Minute,
		"CSRF",
		"X-XSRF-TOKEN")
	userService := usrserv.NewUserService(userStore,
		idGenerator,
		passwordEncrypter,
		passwordValidator)
	documentService := docserv.NewDocumentStore(documentStore, idGenerator)
	collectionService := cllcsserv.NewCollectionService(collectionStore, idGenerator)

	router := http.NewServeMux()
	rest.SetupRestRouter(
		router,
		appConfig,
		metricStore,
		sessionStore,
		userStore,
		sessionService,
		userService,
		documentService,
		collectionService,
	)

	addr := fmt.Sprintf(":%s", appConfig.ServerPort())
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Printf("Server started listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("cmd/main.main: [%v]", err)
		}
	}()

	<-ctx.Done()
	log.Printf("Server stopped listening on %s", addr)
}
