package stores

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"context"
	"io"
	"os"
	"time"
)

type CollectionStore interface {
	Save(ctx context.Context, collection models.Collection) error
	PinDocument(ctx context.Context, collectionId, documentId string) error
	IsExist(ctx context.Context, id string) (bool, error)
	IsPinned(ctx context.Context, collectionId, documentId string) (bool, error)
	Find(ctx context.Context, id string) (*models.Collection, error)
	FindAllByUserId(ctx context.Context, userId string) ([]*models.Collection, error)
	Rename(ctx context.Context, id, newName string) error
	UnpinDocument(ctx context.Context, collectionId, documentId string) error
	Delete(ctx context.Context, id string) error
}

type DocumentStore interface {
	Save(ctx context.Context, document models.Document) error
	IsIdExist(ctx context.Context, id string) (bool, error)
	IsNameExist(ctx context.Context, userId, name string) (bool, error)
	Open(ctx context.Context, id string) (models.Document, error)
	Rename(ctx context.Context, id, newName string) error
	Delete(ctx context.Context, id string) error
}

type FileStore interface {
	Save(filename string, data io.Reader) error
	IsExist(filename string) (bool, error)
	Open(filename string) (*os.File, error)
	Rename(oldName string, newName string) error
	Delete(filename string) error
}

type UserStore interface {
	Register(ctx context.Context, user models.User) error
	IsIdRegistered(ctx context.Context, id string) (bool, error)
	IsNameRegistered(ctx context.Context, name string) (bool, error)
	Find(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, user models.User) error
	Deregister(ctx context.Context, id string) error
}

type MetricStore interface {
	Track(ctx context.Context, metric models.Metric) error
	IsTracked(ctx context.Context, timestamp time.Time, name models.MetricName) (bool, error)
	Find(ctx context.Context, timestamp time.Time, name models.MetricName) (models.Metric, error)
	FindAllByTimestamp(ctx context.Context, timestamp time.Time) ([]models.Metric, error)
	FindAllByName(ctx context.Context, name models.MetricName) ([]models.Metric, error)
}
