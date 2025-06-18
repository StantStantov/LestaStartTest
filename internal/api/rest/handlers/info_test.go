//go:build integration || !unit

package handlers_test

import (
	"Stant/LestaGamesInternship/internal/api/rest/dto"
	"Stant/LestaGamesInternship/internal/api/rest/handlers"
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/infra/pgsql"
	"Stant/LestaGamesInternship/internal/pkg/apptest"
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"
	"testing"
)

func TestRestApi(t *testing.T) {
	t.SkipNow()
	ctx := context.Background()

	dbPool := apptest.GetTestPool(t, ctx, os.Getenv(config.DatabaseUrlEnv))
	tx := apptest.GetTestTx(t, ctx, dbPool)

	config, _ := config.ReadAppConfig()
	metricsStore := pgsql.NewMetricStore(tx)

	router := http.NewServeMux()
	router.Handle("GET /api/status", handlers.HandleGetStatus())
	router.Handle("GET /api/metrics", handlers.HandleGetMetrics(metricsStore))
	router.Handle("GET /api/version", handlers.HandleGetVersion(config))

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
		wantBody := dto.AppMetrics{}

		request, err := http.NewRequest("GET", "/api/metrics", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		gotCode := response.Code
		gotBody := dto.AppMetrics{}
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
