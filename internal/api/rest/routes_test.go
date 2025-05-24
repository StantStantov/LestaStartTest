package rest_test

import (
	"Stant/LestaGamesInternship/internal/api/rest"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestApi(t *testing.T) {
	router := http.NewServeMux()
	router.Handle("GET /api/status", rest.HandleStatusGet())

	t.Run("Get app status", func(t *testing.T) {
		t.Helper()

		wantCode:= http.StatusOK
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
}
