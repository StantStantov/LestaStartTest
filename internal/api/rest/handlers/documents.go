package handlers

import (
	"Stant/LestaGamesInternship/internal/api/rest/dto"
	"Stant/LestaGamesInternship/internal/app/services/cllcsserv"
	"Stant/LestaGamesInternship/internal/app/services/docserv"
	"Stant/LestaGamesInternship/internal/app/services/huffman"
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"Stant/LestaGamesInternship/internal/app/services/tfidf"
	"Stant/LestaGamesInternship/internal/domain/models"
	"cmp"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
)

func HandlePostDocument(
	formFileKey string,
	documentService *docserv.DocumentService,
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

		file, fileHeader, err := r.FormFile(formFileKey)
		if err != nil {
			log.Printf("handlers/documents.HandlePostDocument: [%v]", err)
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		if err := documentService.Upload(r.Context(), *userId, fileHeader.Filename, file); err != nil {
			log.Printf("handlers/documents.HandlePostDocument: [%v]", err)
			http.Error(w, `{"error": "User ID empty"}`, http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func HandleGetDocument(
	pathValue string,
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

		documentId := r.PathValue(pathValue)
		if documentId == "" {
			http.Error(w, `{"error": "Empty path value"}`, http.StatusBadRequest)
			return
		}

		accessable, err := documentService.HasAccess(r.Context(), *userId, documentId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		document, err := documentService.Get(r.Context(), documentId)
		if err != nil {
			log.Printf("handlers/documents.HandleGetDocument: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		builder := new(strings.Builder)
		io.Copy(builder, document.File())

		body := dto.DocumentWithData{
			Document: dto.Document{Id: document.Id(), Name: document.Name()},
			Data:     builder.String(),
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("handlers/documents.HandleGetDocument: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func HandleGetDocuments(
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

		documents, err := documentService.GetAll(r.Context(), *userId)
		if err != nil && errors.Unwrap(errors.Unwrap(err)) != pgx.ErrNoRows {
			log.Printf("handlers/documents.HandleGetDocuments: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		body := make([]dto.Document, len(documents))
		for i, document := range documents {
			body[i] = dto.Document{Id: document.Id(), Name: document.Name()}
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("handlers/documents.HandleGetDocuments: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func HandleGetDocumentStatistics(
	pathValue string,
	documentService *docserv.DocumentService,
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
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		documentId := r.PathValue(pathValue)
		if documentId == "" {
			http.Error(w, `{"error": "Empty Document ID"}`, http.StatusBadRequest)
			return
		}

		accessable, err := documentService.HasAccess(r.Context(), *userId, documentId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		mainDocument, err := documentService.Get(r.Context(), documentId)
		if err != nil {
			log.Printf("handlers/documents.HandleGetDocumentStatistics: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		terms, err := tfidf.ProcessReaderToTerms(mainDocument.File())
		if err != nil {
			log.Printf("handlers/documents.HandleGetDocumentStatistics: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		termsFrequencies := tfidf.GetTermFrequency(terms)

		termsAmount := uint64(0)
		otherCollections, err := collectionService.GetAllByDocumentId(r.Context(), documentId)
		if err != nil {
			log.Printf("handlers/documents.HandleGetDocumentStatistics: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		if len(otherCollections) == 0 {
			http.Error(w, `{"error": "Document isn't pinned to any Collection"}`, http.StatusInternalServerError)
			return
		}

		for _, collection := range otherCollections {
			documents := collection.Documents()
			for _, document := range documents {
				if document.Id() == documentId {
					continue
				}
				terms, err := tfidf.ProcessReaderToTerms(document.File())
				if err != nil {
					log.Printf("handlers/documents.HandleGetDocumentStatistics: [%v]", err)
					http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
					return
				}

				termsAmount += uint64(len(terms))
			}
		}

		body := make([]dto.Term, 0, models.MaxStatisticTermsAmount)
		for term, amount := range termsFrequencies {
			body = append(body, dto.Term{Word: term, Tf: amount, Idf: tfidf.CalculateIdf(termsAmount, amount)})
		}
		slices.SortFunc(body, func(E1, E2 dto.Term) int { return cmp.Compare(E1.Idf, E2.Idf) })

		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("handlers/documents.HandleGetDocumentStatistics: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func HandleGetDocumentHuffman(
	pathValue string,
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
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		documentId := r.PathValue(pathValue)
		if documentId == "" {
			http.Error(w, `{"error": "Empty Document ID"}`, http.StatusBadRequest)
			return
		}

		accessable, err := documentService.HasAccess(r.Context(), *userId, documentId)
		if err != nil || !accessable {
			http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
			return
		}

		document, err := documentService.Get(r.Context(), documentId)
		if err != nil {
			log.Printf("handlers/documents.HandleGetDocumentHuffman: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		encoded, codesTable, err := huffman.Encode(document.File())
		if err != nil {
			log.Printf("handlers/documents.HandleGetDocumentHuffman: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		body := dto.HuffmanEncoding{Codes: codesTable, Data: encoded}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("handlers/documents.HandleGetDocumentHuffman: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func HandleDeleteDocument(
	pathValue string,
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

		documentId := r.PathValue(pathValue)
		if documentId == "" {
			http.Error(w, `{"error": "Empty path value"}`, http.StatusBadRequest)
			return
		}

		if err := documentService.Delete(r.Context(), documentId); err != nil {
			log.Printf("handlers/documents.HandleDeleteDocument: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
