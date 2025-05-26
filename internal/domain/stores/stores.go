package stores

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"time"
)

type TermStore interface {
	crudStore[int, models.Term]
}

type MetricStore interface {
	crudStore[time.Time, models.Metric]
	ReadAllByTimestamp(timestamp time.Time) ([]models.Metric, error)
	ReadAllByName(name models.MetricName) ([]models.Metric, error)
}

type crudStore[K any, V any] interface {
	Create(value V) error
	Read(key K) (V, error)
	ReadAll() ([]V, error)
	CountAll() (int, error)
	Update(key K, value V) error
	Delete(key K) error
}
