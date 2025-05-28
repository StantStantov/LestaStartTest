package rest_test

import (
	"Stant/LestaGamesInternship/internal/api/rest"
	"Stant/LestaGamesInternship/internal/app/services/metricService"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"
)

func TestRestApi(t *testing.T) {
	metricsStore := stores.NewInMemoryMetricStore()

	router := http.NewServeMux()
	rest.SetupRestRouter(router, metricsStore)

	t.Run("Get app status", func(t *testing.T) {
		t.Helper()

		wantCode := http.StatusOK
		wantBody := map[string]string{"status": "OK"}

		request, err := http.NewRequest("GET", "/api/status", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		gotCode := response.Code
		gotBody := map[string]string{}
		json.NewDecoder(response.Body).Decode(&gotBody)

		if gotCode != wantCode {
			t.Fatalf("Want Status %d, got %d", wantCode, gotCode)
		}
		if !maps.Equal(wantBody, gotBody) {
			t.Fatalf("Want Body %v, got %v", wantBody, gotBody)
		}
	})
	t.Run("Get app empty metrics", func(t *testing.T) {
		t.Helper()

		wantCode := http.StatusOK
		wantBody := metricService.AppMetrics{}

		request, err := http.NewRequest("GET", "/api/metrics", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		gotCode := response.Code
		gotBody := metricService.AppMetrics{}
		json.NewDecoder(response.Body).Decode(&gotBody)

		if gotCode != wantCode {
			t.Fatalf("Want Status %d, got %d", wantCode, gotCode)
		}
		if wantBody != gotBody {
			t.Fatalf("Want Body %+v, got %+v", wantBody, gotBody)
		}
	})
	t.Run("Get app version", func(t *testing.T) {
		t.Helper()

		wantCode := http.StatusOK
		info, _ := debug.ReadBuildInfo()
		wantBody := map[string]string{"version": info.Main.Version}

		request, err := http.NewRequest("GET", "/api/version", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		gotCode := response.Code
		gotBody := map[string]string{}
		json.NewDecoder(response.Body).Decode(&gotBody)

		if gotCode != wantCode {
			t.Fatalf("Want Status %d, got %d", wantCode, gotCode)
		}
		if !maps.Equal(wantBody, gotBody) {
			t.Fatalf("Want Body %v, got %v", wantBody, gotBody)
		}
	})
}
