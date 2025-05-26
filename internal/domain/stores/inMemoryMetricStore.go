package stores

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"fmt"
	"slices"
	"time"
)

type InMemoryMetricStore struct {
	values []models.Metric
	length int
}

func NewInMemoryMetricStore() *InMemoryMetricStore {
	return &InMemoryMetricStore{values: []models.Metric{}, length: 0}
}

func (ts *InMemoryMetricStore) Create(metric models.Metric) error {
	ts.values = append(ts.values, metric)
	ts.length++
	return nil
}

func (ts *InMemoryMetricStore) Read(key time.Time) (models.Metric, error) {
	index := slices.IndexFunc(ts.values, isEqualByTimestamp(key))
	if index == -1 {
		return models.Metric{}, fmt.Errorf("InMemoryTermStore.Read: Metric does not exist")
	}
	return ts.values[index], nil
}

func (ts *InMemoryMetricStore) ReadAll() ([]models.Metric, error) {
	return slices.Clone(ts.values), nil
}

func (ts *InMemoryMetricStore) ReadAllByTimestamp(timestamp time.Time) ([]models.Metric, error) {
	metrics := []models.Metric{}
	for _, metric := range ts.values {
		if metric.Timestamp().Equal(timestamp) {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func (ts *InMemoryMetricStore) ReadAllByName(name models.MetricName) ([]models.Metric, error) {
	metrics := []models.Metric{}
	for _, metric := range ts.values {
		if metric.Name() == name {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func (ts *InMemoryMetricStore) CountAll() (int, error) {
	return ts.length, nil
}

func (ts *InMemoryMetricStore) Update(key time.Time, metric models.Metric) error {
	index := slices.IndexFunc(ts.values, isEqualByTimestamp(key))
	if index == -1 {
		return fmt.Errorf("InMemoryMetricStore.Update: Metric does not exist")
	}
	ts.values[index] = models.NewMetric(key, models.MetricName(metric.Name()), metric.Value())
	return nil
}

func (ts *InMemoryMetricStore) Delete(key time.Time) error {
	index := slices.IndexFunc(ts.values, isEqualByTimestamp(key))
	if index == -1 {
		return fmt.Errorf("InMemoryMetricStore.Delete: Metric does not exist")
	}
	ts.values = slices.DeleteFunc(ts.values, isEqualByTimestamp(key))
	ts.length--
	return nil
}

func isEqualByTimestamp(T time.Time) func(models.Metric) bool {
	return func(E models.Metric) bool {
		return E.Timestamp().Equal(T)
	}
}
