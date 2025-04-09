package internal

import (
	"Stant/LestaGamesInternship/internal/views"
	"log"
	"net/http"
)

func HandleIndexGet() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})
}

func HandleIndexPost() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		_, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			return
		}

		println(header.Filename)
	})
}
