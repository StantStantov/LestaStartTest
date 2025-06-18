//go:build unit || !integration

package mtrcserv_test

import (
	"Stant/LestaGamesInternship/internal/app/services/mtrcserv"
	"Stant/LestaGamesInternship/internal/domain/models"
	"testing"
	"time"
)

func TestMetricService(t *testing.T) {
	t.Run("Test SumValues", func(t *testing.T) {
		testSumValues(t)
	})
	t.Run("Test FindMaxByValue", func(t *testing.T) {
		testFindMaxByValue(t)
	})
	t.Run("Test FindMaxByTimestamp", func(t *testing.T) {
		testFindMaxByTimestamp(t)
	})
	t.Run("Test FindMinByValue", func(t *testing.T) {
		testFindMinByValue(t)
	})
	t.Run("Test FindMinByTimestamp", func(t *testing.T) {
		testFindMinByTimestamp(t)
	})
}

func testSumValues(t *testing.T) {
	t.Helper()

	timestamp := time.Date(2025, time.May, 26, 20, 0, 0, 0, time.UTC)
	metrics := []models.Metric{
		models.NewMetric(timestamp, models.FilesProcessed, 1),
		models.NewMetric(timestamp, models.FilesProcessed, 1),
		models.NewMetric(timestamp, models.FilesProcessed, 1),
		models.NewMetric(timestamp, models.FilesProcessed, 1),
	}

	want := 4.0
	got := mtrcserv.SumValues(metrics)
	if got != want {
		t.Errorf("Wanted %f, got %f", want, got)
	}
}

func testFindMaxByValue(t *testing.T) {
	t.Helper()

	timestamp := time.Date(2025, time.May, 26, 20, 0, 0, 0, time.UTC)
	want := 13.49

	metrics := []models.Metric{
		models.NewMetric(timestamp, models.TimeProcessed, 3.6),
		models.NewMetric(timestamp, models.TimeProcessed, 11.1),
		models.NewMetric(timestamp, models.TimeProcessed, want),
		models.NewMetric(timestamp, models.TimeProcessed, 5.9),
	}
	got, err := mtrcserv.FindMaxByValue(metrics)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Wanted %f, got %f", want, got)
	}
}

func testFindMaxByTimestamp(t *testing.T) {
	t.Helper()

	timestamp := time.Date(2025, time.May, 26, 20, 0, 0, 0, time.UTC)
	want := timestamp.AddDate(0, 0, 1)

	metrics := []models.Metric{
		models.NewMetric(want, models.TimeProcessed, 3.6),
		models.NewMetric(timestamp, models.TimeProcessed, 11.1),
		models.NewMetric(timestamp, models.TimeProcessed, 13.2),
		models.NewMetric(timestamp, models.TimeProcessed, 5.9),
	}
	got, err := mtrcserv.FindMaxByTimestamp(metrics)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(want) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func testFindMinByValue(t *testing.T) {
	t.Helper()

	timestamp := time.Date(2025, time.May, 26, 20, 0, 0, 0, time.UTC)
	want := 0.49

	metrics := []models.Metric{
		models.NewMetric(timestamp, models.TimeProcessed, 11.1),
		models.NewMetric(timestamp, models.TimeProcessed, want),
		models.NewMetric(timestamp, models.TimeProcessed, 3.6),
		models.NewMetric(timestamp, models.TimeProcessed, 5.9),
	}
	got, err := mtrcserv.FindMinByValue(metrics)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Wanted %f, got %f", want, got)
	}
}

func testFindMinByTimestamp(t *testing.T) {
	t.Helper()

	timestamp := time.Date(2025, time.May, 26, 20, 0, 0, 0, time.UTC)
	want := timestamp.AddDate(0, 0, -1)

	metrics := []models.Metric{
		models.NewMetric(timestamp, models.TimeProcessed, 11.1),
		models.NewMetric(timestamp, models.TimeProcessed, 43),
		models.NewMetric(want, models.TimeProcessed, 3.6),
		models.NewMetric(timestamp, models.TimeProcessed, 5.9),
	}
	got, err := mtrcserv.FindMinByTimestamp(metrics)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(want) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}
