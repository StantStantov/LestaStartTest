//go:build integration || !unit

package middlewares_test

import (
	"Stant/LestaGamesInternship/internal/api/rest/middlewares"
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"Stant/LestaGamesInternship/internal/pkg/apptest"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	ctx := context.Background()

	dbPool := apptest.GetTestPool(t, ctx, os.Getenv(config.DatabaseUrlEnv))

	t.Run("Collect if OK", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		metricStore := pgsql.NewMetricStore(tx)

		wantLength := 1

		router := newMockRouter(metricStore)
		request, err := http.NewRequest(http.MethodPost, "/good", nil)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		router.ServeHTTP(httptest.NewRecorder(), request)

		metrics, err := metricStore.FindAllByName(ctx, models.FilesProcessed)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		gotLength := len(metrics)

		if wantLength != gotLength {
			t.Errorf("Wanted %d, got %d", wantLength, gotLength)
		}
	})
	t.Run("Do nothing if Error", func(t *testing.T) {
		t.Parallel()

		tx := apptest.GetTestTx(t, ctx, dbPool)
		metricStore := pgsql.NewMetricStore(tx)

		router := newMockRouter(metricStore)
		request, err := http.NewRequest(http.MethodPost, "/bad", nil)
		if err != nil {
			t.Fatalf("Wanted %v, got %v", nil, err)
		}
		router.ServeHTTP(httptest.NewRecorder(), request)

		if _, err := metricStore.FindAllByName(ctx, models.FilesProcessed); err == nil {
			t.Errorf("Wanted err, got %v", err)
		}
	})
}

func newMockRouter(metricStore stores.MetricStore) *http.ServeMux {
	metricMiddleware := middlewares.NewCollectMetricsMiddleware(metricStore)

	router := http.NewServeMux()
	router.Handle("POST /good", metricMiddleware(mockPostHandler(http.StatusOK)))
	router.Handle("POST /bad", metricMiddleware(mockPostHandler(http.StatusInternalServerError)))

	return router
}

func mockPostHandler(code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
}
