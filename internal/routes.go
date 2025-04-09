package internal

import (
	"Stant/LestaGamesInternship/internal/views"
	"bufio"
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

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			return
		}
	})
}
