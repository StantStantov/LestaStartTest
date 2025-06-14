//go:build integration || !unit

package pgsql_test

import (
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"Stant/LestaGamesInternship/internal/pkg/apptest"
	"context"
	"os"
	"slices"
	"testing"
	"time"
)

func TestMetricStore(t *testing.T) {
	ctx := context.Background()

	dbPool := apptest.GetTestPool(t, ctx, os.Getenv(config.DatabaseUrlEnv))

	t.Run("Test Track Metric", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		metricStore := pgsql.NewMetricStore(tx)

		testMetricStoreTrack(t, ctx, metricStore)
	})
	t.Run("Test Read Metric", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		metricStore := pgsql.NewMetricStore(tx)

		testMetricStoreFind(t, ctx, metricStore)
	})
	t.Run("Test ReadAll Metrics", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		metricStore := pgsql.NewMetricStore(tx)

		testMetricStoreFindAll(t, ctx, metricStore)
	})
}

func testMetricStoreTrack(t *testing.T, ctx context.Context, metricStore stores.MetricStore) {
	t.Helper()

	t.Run("PASS if tracked", func(t *testing.T) {
		want := true
		metric := models.NewMetric(time.Now(), models.TimeProcessed, 13.5)

		if err := metricStore.Track(ctx, metric); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got, err := metricStore.IsTracked(ctx, metric.Timestamp(), metric.Name())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != got {
			t.Fatalf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("PASS if updated", func(t *testing.T) {
		want := true
		metric := models.NewMetric(time.Now(), models.TimeProcessed, 13.5)
		anotherMetric := models.NewMetric(metric.Timestamp(), metric.Name(), 21.32)

		if err := metricStore.Track(ctx, metric); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		if err := metricStore.Track(ctx, anotherMetric); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got, err := metricStore.IsTracked(ctx, anotherMetric.Timestamp(), anotherMetric.Name())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if want != got {
			t.Fatalf("Wanted %v, got %v", want, got)
		}
	})
}

func testMetricStoreFind(t *testing.T, ctx context.Context, metricStore stores.MetricStore) {
	t.Helper()

	t.Run("PASS if found", func(t *testing.T) {
		wantTime := time.Now().Round(time.Second)
		wantName := models.FilesProcessed
		wantValue := 1.0
		want := models.NewMetric(wantTime, wantName, wantValue)

		if err := metricStore.Track(ctx, want); err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		got, err := metricStore.Find(ctx, want.Timestamp(), want.Name())
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotTime := got.Timestamp().Round(time.Second)
		gotName := got.Name()
		gotValue := got.Value()

		if !wantTime.Equal(gotTime) {
			t.Errorf("Wanted %v, got %v", wantTime, gotTime)
		}
		if wantName != gotName {
			t.Errorf("Wanted %v, got %v", wantName, gotName)
		}
		if wantValue != gotValue {
			t.Errorf("Wanted %v, got %v", wantValue, gotValue)
		}
	})
	t.Run("FAIL if doesn't exist", func(t *testing.T) {
		if _, err := metricStore.Find(ctx, time.Now(), models.FilesProcessed); err == nil {
			t.Fatalf("Wanted err, got %v", err)
		}
	})
}

func testMetricStoreFindAll(t *testing.T, ctx context.Context, metricStore stores.MetricStore) {
	t.Helper()

	isEqualMetrics := func(E1, E2 models.Metric) bool {
		isEqualTime := E1.Timestamp().Round(time.Second).Equal(E2.Timestamp().Round(time.Second))
		isEqualName := E1.Name() == E2.Name()
		isEqualValue := E1.Value() == E2.Value()
		return isEqualTime && isEqualName && isEqualValue
	}

	t.Run("PASS if found all by Timestamp", func(t *testing.T) {
		wantTime := time.Now().Round(time.Second)
		want := []models.Metric{
			models.NewMetric(wantTime, models.FilesProcessed, 1.0),
			models.NewMetric(wantTime, models.TimeProcessed, 5.1),
		}

		for _, metric := range want {
			if err := metricStore.Track(ctx, metric); err != nil {
				t.Fatalf("Wanted %v, got %v", nil, err)
			}
		}
		got, err := metricStore.FindAllByTimestamp(ctx, wantTime)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if !slices.EqualFunc(want, got, isEqualMetrics) {
			t.Fatalf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("PASS if found all by Name", func(t *testing.T) {
		wantName := models.TimeProcessed
		want := []models.Metric{
			models.NewMetric(time.Now().Round(time.Second), wantName, 3.0),
			models.NewMetric(time.Now().Add(15*time.Second).Round(time.Second), wantName, 5.1),
		}

		for _, metric := range want {
			if err := metricStore.Track(ctx, metric); err != nil {
				t.Fatalf("Wanted %v, got %v", nil, err)
			}
		}
		got, err := metricStore.FindAllByName(ctx, wantName)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}

		if !slices.EqualFunc(want, got, isEqualMetrics) {
			t.Fatalf("Wanted %v, got %v", want, got)
		}
	})
	t.Run("FAIL if Timestamp doesn't exist", func(t *testing.T) {
		_, err := metricStore.FindAllByTimestamp(ctx, time.Time{})
		if err == nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
	})
	t.Run("FAIL if Name doesn't exist", func(t *testing.T) {
		_, err := metricStore.FindAllByName(ctx, "")
		if err == nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
	})
}
