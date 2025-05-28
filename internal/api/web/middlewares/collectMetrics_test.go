package middlewares_test

import (
	"Stant/LestaGamesInternship/internal/api/web/middlewares"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	// metricStore := stores.NewInMemoryMetricStore()
	//
	// metricMiddleware := middlewares.NewCollectMetricsMiddleware(metricStore)
	//
	// router := http.NewServeMux()
	// router.Handle("POST /good", metricMiddleware(mockPostHandler(map[string]string{"Content-Type": "multipart/form-data"})))
	// router.Handle("POST /bad", metricMiddleware(mockPostHandler(nil)))

	t.Run("Collect if OK", func(t *testing.T) {
		t.Helper()

		wantLength := 2

		metricStore := stores.NewInMemoryMetricStore()
		router := newMockRouter(metricStore)
		request, err := http.NewRequest(http.MethodPost, "/good", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(httptest.NewRecorder(), request)

		metrics, err := metricStore.ReadAll()
		if err != nil {
			t.Fatal(err)
		}
		gotLength := len(metrics)

		if wantLength != gotLength {
			t.Errorf("Wanted %d, got %d", wantLength, gotLength)
		}
	})
	t.Run("Do nothing if Error", func(t *testing.T) {
		t.Helper()

		wantLength := 0

		metricStore := stores.NewInMemoryMetricStore()
		router := newMockRouter(metricStore)
		request, err := http.NewRequest(http.MethodPost, "/bad", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(httptest.NewRecorder(), request)

		metrics, err := metricStore.ReadAll()
		if err != nil {
			t.Fatal(err)
		}
		gotLength := len(metrics)

		if wantLength != gotLength {
			t.Errorf("Wanted %d, got %d", wantLength, gotLength)
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
