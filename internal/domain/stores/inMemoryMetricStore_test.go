package stores_test

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"slices"
	"testing"
	"time"
)

func TestInMemoryMetricStore(t *testing.T) {
	t.Run("Test Create", func(t *testing.T) {
		metricStore := stores.NewInMemoryMetricStore()
		testMetricStoreCreate(metricStore, t)
	})
	t.Run("Test Read", func(t *testing.T) {
		metricStore := stores.NewInMemoryMetricStore()
		testMetricStoreRead(metricStore, t)
	})
	t.Run("Test ReadAll", func(t *testing.T) {
		metricStore := stores.NewInMemoryMetricStore()
		testMetricStoreReadAll(metricStore, t)
	})
	t.Run("Test ReadAllByTimestamp", func(t *testing.T) {
		metricStore := stores.NewInMemoryMetricStore()
		testMetricStoreReadAllByTimestamp(metricStore, t)
	})
	t.Run("Test ReadAllByName", func(t *testing.T) {
		metricStore := stores.NewInMemoryMetricStore()
		testMetricStoreReadAllByName(metricStore, t)
	})
	t.Run("Test Update", func(t *testing.T) {
		metricStore := stores.NewInMemoryMetricStore()
		testMetricStoreUpdate(metricStore, t)
	})
	t.Run("Test Delete", func(t *testing.T) {
		metricStore := stores.NewInMemoryMetricStore()
		testMetricStoreDelete(metricStore, t)
	})
}

func testMetricStoreCreate(metricStore stores.MetricStore, t *testing.T) {
	t.Helper()

	metric := models.NewMetric(time.Now(), models.FilesProcessed, 1)

	if err := metricStore.Create(metric); err != nil {
		t.Fatal(err)
	}
}

func testMetricStoreRead(metricStore stores.MetricStore, t *testing.T) {
	t.Helper()

	want := models.NewMetric(time.Now(), models.TimeProcessed, 5.1)

	metricStore.Create(want)
	got, err := metricStore.Read(want.Timestamp())
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}

func testMetricStoreReadAll(metricStore stores.MetricStore, t *testing.T) {
	t.Helper()

	timestamp := time.Now()
	want := []models.Metric{
		models.NewMetric(timestamp, models.FilesProcessed, 1),
		models.NewMetric(timestamp, models.TimeProcessed, 32.9),
	}

	for _, metric := range want {
		metricStore.Create(metric)
	}
	got, err := metricStore.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func testMetricStoreReadAllByTimestamp(metricStore stores.MetricStore, t *testing.T) {
	t.Helper()

	wantTimestamp := time.Now()
	anotherTimestamp := time.Now().AddDate(0, 0, 1)
	metrics := []models.Metric{
		models.NewMetric(wantTimestamp, models.FilesProcessed, 1),
		models.NewMetric(wantTimestamp, models.TimeProcessed, 32.9),
		models.NewMetric(anotherTimestamp, models.FilesProcessed, 1),
		models.NewMetric(anotherTimestamp, models.TimeProcessed, 12.7),
	}
	want := []models.Metric{
		metrics[0], metrics[1],
	}

	for _, metric := range metrics {
		metricStore.Create(metric)
	}
	got, err := metricStore.ReadAllByTimestamp(wantTimestamp)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func testMetricStoreReadAllByName(metricStore stores.MetricStore, t *testing.T) {
	t.Helper()

	timestamp := time.Now()
	anotherTimestamp := time.Now().AddDate(0, 0, 1)
	metrics := []models.Metric{
		models.NewMetric(timestamp, models.FilesProcessed, 1),
		models.NewMetric(timestamp, models.TimeProcessed, 32.9),
		models.NewMetric(anotherTimestamp, models.FilesProcessed, 1),
		models.NewMetric(anotherTimestamp, models.TimeProcessed, 12.7),
	}
	want := []models.Metric{
		metrics[0], metrics[2],
	}

	for _, metric := range metrics {
		metricStore.Create(metric)
	}
	got, err := metricStore.ReadAllByName(models.FilesProcessed)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func testMetricStoreUpdate(metricStore stores.MetricStore, t *testing.T) {
	t.Helper()

	oldMetric := models.NewMetric(time.Now(), models.TimeProcessed, 2.1)
	newMetric := models.NewMetric(time.Now(), models.TimeProcessed, 5.9)
	want := models.NewMetric(oldMetric.Timestamp(), newMetric.Name(), newMetric.Value())

	metricStore.Create(oldMetric)
	if err := metricStore.Update(oldMetric.Timestamp(), newMetric); err != nil {
		t.Fatal(err)
	}

	got, err := metricStore.Read(want.Timestamp())
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}

func testMetricStoreDelete(metricStore stores.MetricStore, t *testing.T) {
	t.Helper()

	metric := models.NewMetric(time.Now(), models.TimeProcessed, 2.1)
	want := models.Metric{}

	metricStore.Create(metric)
	if err := metricStore.Delete(metric.Timestamp()); err != nil {
		t.Fatal(err)
	}

	got, err := metricStore.Read(metric.Timestamp())
	if err == nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Wanted %+v, got %+v", want, got)
	}
}
