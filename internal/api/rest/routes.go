package rest

import (
	"fmt"
	"net/http"
)

func HandleStatusGet() http.HandlerFunc {
	status := `{"status": "OK"}`
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, status)
	})
}
