package cllcsserv

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/services"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CollectionService struct {
	collectionStore stores.CollectionStore

	idGen services.IdGenerator
}

func NewCollectionService(
	collectionStore stores.CollectionStore,
	idGen services.IdGenerator,
) *CollectionService {
	return &CollectionService{collectionStore: collectionStore, idGen: idGen}
}

func (s *CollectionService) Create(ctx context.Context, userId, name string) error {
	collection := models.NewEmptyCollection(s.idGen.GenerateId(), userId, name)

	if err := s.collectionStore.Save(ctx, *collection); err != nil {
		return fmt.Errorf("collections/collectionService.Create: [%w]", err)
	}

	return nil
}

func (s *CollectionService) Get(ctx context.Context, collectionId string) (*models.Collection, error) {
	collection, err := s.collectionStore.Find(ctx, collectionId)
	if err != nil {
		return nil, fmt.Errorf("collections/collectionService.Get: [%w]", err)
	}

	return collection, nil
}

func (s *CollectionService) GetAllByUserId(ctx context.Context, userId string) ([]*models.Collection, error) {
	collections, err := s.collectionStore.FindAllByUserId(ctx, userId)
	if err != nil && errors.Unwrap(errors.Unwrap(err)) != pgx.ErrNoRows {
		return nil, fmt.Errorf("collections/collectionService.GetAllByUserId: [%w]", err)
	}

	return collections, nil
}

func (s *CollectionService) GetAllByDocumentId(ctx context.Context, documentId string) ([]*models.Collection, error) {
	collections, err := s.collectionStore.FindAllByDocumentId(ctx, documentId)
	if err != nil && errors.Unwrap(errors.Unwrap(err)) != pgx.ErrNoRows {
		return nil, fmt.Errorf("collections/collectionService.GetAllByDocumentId: [%w]", err)
	}

	return collections, nil
}

func (s *CollectionService) PinToCollection(ctx context.Context, collectionId, documentId string) error {
	isPinned, err := s.collectionStore.IsPinned(ctx, collectionId, documentId)
	if err != nil {
		return fmt.Errorf("collections/collectionService.PinToCollection: [%w]", err)
	}
	if isPinned {
		return nil
	}

	if err := s.collectionStore.PinDocument(ctx, collectionId, documentId); err != nil {
		return fmt.Errorf("collections/collectionService.PinToCollection: [%w]", err)
	}

	return nil
}

func (s *CollectionService) UnpinFromCollection(ctx context.Context, collectionId, documentId string) error {
	isPinned, err := s.collectionStore.IsPinned(ctx, collectionId, documentId)
	if err != nil {
		return fmt.Errorf("collections/collectionService.UnpinFromCollection: [%w]", err)
	}
	if !isPinned {
		return nil
	}

	if err := s.collectionStore.UnpinDocument(ctx, collectionId, documentId); err != nil {
		return fmt.Errorf("collections/collectionService.UnpinFromCollection: [%w]", err)
	}

	return nil
}

func (s *CollectionService) HasAccess(ctx context.Context, userId, collectionId string) (bool, error) {
	owned, err := s.collectionStore.IsOwned(ctx, userId, collectionId)
	if err != nil {
		return false, fmt.Errorf("collections/collectionService.HasAccess: [%w]", err)
	}

	return owned, nil
}
