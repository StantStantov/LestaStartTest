package rest

import (
	"Stant/LestaGamesInternship/internal/api/rest/handlers"
	"Stant/LestaGamesInternship/internal/api/rest/middlewares"
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/app/services/cllcsserv"
	"Stant/LestaGamesInternship/internal/app/services/docserv"
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"Stant/LestaGamesInternship/internal/app/services/usrserv"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"net/http"
)

// Дальше Бога нет...

func SetupRestRouter(
	router *http.ServeMux,
	appConfig *config.AppConfig,
	metricsStore stores.MetricStore,
	sessionStore stores.SessionStore,
	userStore stores.UserStore,
	sessionService *sesserv.SessionService,
	userService *usrserv.UserService,
	documentService *docserv.DocumentService,
	collectionService *cllcsserv.CollectionService,
) {
	checkSession := middlewares.CheckSession(sessionService)
	requireAuth := middlewares.RequireAuth(userStore)
	requireCsrf := middlewares.RequireCsrf(sessionService)

	router.Handle("GET /api/status", handlers.HandleGetStatus())
	router.Handle("GET /api/metrics", handlers.HandleGetMetrics(metricsStore))
	router.Handle("GET /api/version", handlers.HandleGetVersion(appConfig))

	router.Handle("POST /api/register",
		checkSession(
			handlers.HandlePostRegister(userService)))
	router.Handle("POST /api/login",
		checkSession(
			handlers.HandlePostLogin(userService, userStore, sessionService)))
	router.Handle("GET /api/logout",
		checkSession(
			requireAuth(
				handlers.HandlePostLogout(sessionService))))
	router.Handle("PATCH /api/user/{user_id}",
		checkSession(
			requireAuth(
				requireCsrf(
					handlers.HandlePatchUser("user_id", userService)))))
	router.Handle("DELETE /api/user/{user_id}",
		checkSession(
			requireAuth(
				requireCsrf(
					handlers.HandleDeleteUser("user_id", userService, sessionService)))))

	router.Handle("POST /api/documents/",
		checkSession(
			requireAuth(
				requireCsrf(
					handlers.HandlePostDocument("document", documentService)))))
	router.Handle("GET /api/documents/{document_id}",
		checkSession(
			requireAuth(
				handlers.HandleGetDocument("document_id", documentService))))
	router.Handle("GET /api/documents/{document_id}/statistics",
		checkSession(
			requireAuth(
				handlers.HandleGetDocumentStatistics("document_id", documentService, collectionService))))
	router.Handle("GET /api/documents/{document_id}/huffman",
		checkSession(
			requireAuth(
				handlers.HandleGetDocumentHuffman("document_id", documentService))))
	router.Handle("GET /api/documents/",
		checkSession(
			requireAuth(
				handlers.HandleGetDocuments(documentService))))
	router.Handle("DELETE /api/documents/{document_id}",
		checkSession(
			requireAuth(
				requireCsrf(
					handlers.HandleDeleteDocument("document_id", documentService)))))

	router.Handle("POST /api/collections/",
		checkSession(
			requireAuth(
				requireCsrf(
					handlers.HandlePostCollection(collectionService)))))
	router.Handle("GET /api/collections/",
		checkSession(
			requireAuth(
				handlers.HandleGetCollections(collectionService))))
	router.Handle("GET /api/collections/{collection_id}",
		checkSession(
			requireAuth(
				handlers.HandleGetCollection("collection_id", collectionService))))
	router.Handle("GET /api/collections/{collection_id}/statistics",
		checkSession(
			requireAuth(
				handlers.HandleGetCollectionStatistics("collection_id", collectionService))))
	router.Handle("POST /api/collection/{collection_id}/{document_id}",
		checkSession(
			requireAuth(
				requireCsrf(
					handlers.HandlePostDocumentToCollection("collection_id", "document_id", collectionService, documentService)))))
	router.Handle("DELETE /api/collection/{collection_id}/{document_id}",
		checkSession(
			requireAuth(
				requireCsrf(
					handlers.HandleDeleteDocumentFromCollection("collection_id", "document_id", collectionService, documentService)))))
}
