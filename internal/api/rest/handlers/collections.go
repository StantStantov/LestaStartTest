package handlers

import (
	"Stant/LestaGamesInternship/internal/api/rest/dto"
	"Stant/LestaGamesInternship/internal/app/services/cllcsserv"
	"Stant/LestaGamesInternship/internal/app/services/docserv"
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"Stant/LestaGamesInternship/internal/app/services/tfidf"
	"Stant/LestaGamesInternship/internal/domain/models"
	"cmp"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5"
)

func HandlePostCollection(
	collectionService *cllcsserv.CollectionService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error" : "Session does not exist"}`, http.StatusBadRequest)
			return
		}
		userId := session.UserId()
		if userId == nil {
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		body := struct {
			Name string `json:"name"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("handlers/collections.HandlePostCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		if err := collectionService.Create(r.Context(), *userId, body.Name); err != nil {
			log.Printf("handlers/collections.HandlePostCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func HandleGetCollection(
	pathValue string,
	collectionService *cllcsserv.CollectionService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error": "Session does not exist"}`, http.StatusBadRequest)
			return
		}
		userId := session.UserId()
		if userId == nil {
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		collectionId := r.PathValue(pathValue)
		if collectionId == "" {
			http.Error(w, `{"error": "Empty path value"}`, http.StatusBadRequest)
			return
		}

		accessable, err := collectionService.HasAccess(r.Context(), *userId, collectionId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		collection, err := collectionService.Get(r.Context(), collectionId)
		if err != nil && errors.Unwrap(errors.Unwrap(err)) != pgx.ErrNoRows {
			log.Printf("handlers/collections.HandleGetCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		documents := collection.Documents()

		body := make([]dto.Document, len(documents))
		for i, document := range documents {
			body[i] = dto.Document{Id: document.Id(), Name: document.Name()}
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("handlers/collections.HandleGetCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func HandleGetCollectionStatistics(
	pathValue string,
	collectionService *cllcsserv.CollectionService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error": "Session does not exist"}`, http.StatusBadRequest)
			return
		}
		userId := session.UserId()
		if userId == nil {
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		collectionId := r.PathValue(pathValue)
		if collectionId == "" {
			http.Error(w, `{"error": "Empty Collection ID"}`, http.StatusBadRequest)
			return
		}

		accessable, err := collectionService.HasAccess(r.Context(), *userId, collectionId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		collection, err := collectionService.Get(r.Context(), collectionId)
		if err != nil {
			log.Printf("handlers/collections.HandleGetCollectionStatistics: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		if collection.UserId() != *userId {
			http.Error(w, `{"error": "Not authorized"}`, http.StatusUnauthorized)
			return
		}

		documents := collection.Documents()
		if len(documents) == 0 {
			http.Error(w, `{"error": "Collection doesn't contain any document"}`, http.StatusBadRequest)
			return
		}

		allTerms := []string{}
		for _, document := range documents {
			terms, err := tfidf.ProcessReaderToTerms(document.File())
			if err != nil {
				log.Printf("handlers/collections.HandleGetCollectionStatistics: [%v]", err)
				http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
				return
			}

			allTerms = append(allTerms, terms...)
		}
		termsAmount := uint64(len(allTerms))
		termsFrequencies := tfidf.GetTermFrequency(allTerms)

		body := make([]dto.Term, 0, models.MaxStatisticTermsAmount)
		for term, amount := range termsFrequencies {
			body = append(body, dto.Term{Word: term, Tf: amount, Idf: tfidf.CalculateIdf(termsAmount, amount)})
		}
		slices.SortFunc(body, func(E1, E2 dto.Term) int { return cmp.Compare(E1.Idf, E2.Idf) })

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("handlers/collections.HandleGetCollectionStatistics: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func HandleGetCollections(
	collectionService *cllcsserv.CollectionService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error": "Session does not exist"}`, http.StatusBadRequest)
			return
		}
		userId := session.UserId()
		if userId == nil {
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		collections, err := collectionService.GetAllByUserId(r.Context(), *userId)
		if err != nil {
			log.Printf("handlers/collections.HandleGetCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		body := make([]dto.Collection, len(collections))
		for i, collection := range collections {
			documents := collection.Documents()

			documentsDto := make([]dto.Document, len(documents))
			for i, document := range documents {
				documentsDto[i] = dto.Document{Id: document.Id(), Name: document.Name()}
			}

			body[i] = dto.Collection{
				Id:        collection.Id(),
				Name:      collection.Name(),
				Documents: documentsDto,
			}
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("handlers/collections.HandleGetCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func HandlePostDocumentToCollection(
	collectionPathValue string,
	documentPathValue string,
	collectionService *cllcsserv.CollectionService,
	documentService *docserv.DocumentService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error": "Session does not exist"}`, http.StatusBadRequest)
			return
		}
		userId := session.UserId()
		if userId == nil {
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		collectionId := r.PathValue(collectionPathValue)
		if collectionId == "" {
			http.Error(w, `{"error": "Empty path value"}`, http.StatusBadRequest)
			return
		}
		documentId := r.PathValue(documentPathValue)
		if documentId == "" {
			http.Error(w, `{"error": "Empty path value"}`, http.StatusBadRequest)
			return
		}

		accessable, err := collectionService.HasAccess(r.Context(), *userId, collectionId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}
		accessable, err = documentService.HasAccess(r.Context(), *userId, documentId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		if err := collectionService.PinToCollection(r.Context(), collectionId, documentId); err != nil {
			log.Printf("handlers/collections.HandlePostDocumentToCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func HandleDeleteDocumentFromCollection(
	collectionPathValue string,
	documentPathValue string,
	collectionService *cllcsserv.CollectionService,
	documentService *docserv.DocumentService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error": "Session does not exist"}`, http.StatusBadRequest)
			return
		}
		userId := session.UserId()
		if userId == nil {
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		collectionId := r.PathValue(collectionPathValue)
		if collectionId == "" {
			http.Error(w, `{"error": "Empty path value"}`, http.StatusBadRequest)
			return
		}
		documentId := r.PathValue(documentPathValue)
		if documentId == "" {
			http.Error(w, `{"error": "Empty path value"}`, http.StatusBadRequest)
			return
		}

		accessable, err := collectionService.HasAccess(r.Context(), *userId, collectionId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}
		accessable, err = documentService.HasAccess(r.Context(), *userId, documentId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		if err := collectionService.UnpinFromCollection(r.Context(), collectionId, documentId); err != nil {
			log.Printf("handlers/collections.HandleDeleteDocumentFromCollection: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
