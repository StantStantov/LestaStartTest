package docserv

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/services"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"
)

type DocumentService struct {
	documentStore stores.DocumentStore

	idGen services.IdGenerator
}

func NewDocumentStore(
	documentStore stores.DocumentStore,
	idGen services.IdGenerator,
) *DocumentService {
	return &DocumentService{documentStore: documentStore, idGen: idGen}
}

func (s *DocumentService) Upload(ctx context.Context, userId, filename string, data io.Reader) error {
	document := models.NewDocument(s.idGen.GenerateId(), userId, filename, data)

	if err := s.documentStore.Save(ctx, document); err != nil {
		return fmt.Errorf("documents/documentService.Upload: [%w]", err)
	}

	return nil
}

func (s *DocumentService) Get(ctx context.Context, documentId string) (models.Document, error) {
	document, err := s.documentStore.Open(ctx, documentId)
	if err != nil {
		return models.Document{}, fmt.Errorf("documents/documentService.Get: [%w]", err)
	}

	return document, nil
}

func (s *DocumentService) GetAll(ctx context.Context, userId string) ([]models.Document, error) {
	documents, err := s.documentStore.OpenAll(ctx, userId)
	if err != nil && errors.Unwrap(errors.Unwrap(err)) != pgx.ErrNoRows {
		return nil, fmt.Errorf("documents/documentService.GetAll: [%w]", err)
	}

	return documents, nil
}

func (s *DocumentService) Delete(ctx context.Context, documentId string) error {
	if err := s.documentStore.Delete(ctx, documentId); err != nil {
		return fmt.Errorf("documents/documentService.Delete: [%w]", err)
	}

	return nil
}

func (s *DocumentService) HasAccess(ctx context.Context, userId, documentId string) (bool, error) {
	owned, err := s.documentStore.IsOwned(ctx, userId, documentId)
	if err != nil {
		return false, fmt.Errorf("documents/documentService.HasAccess: [%w]", err)
	}

	return owned, nil
}
