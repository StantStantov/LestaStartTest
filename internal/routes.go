package internal

import (
	"Stant/LestaGamesInternship/internal/views"
	"net/http"
)

func HandleIndex() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})
}
