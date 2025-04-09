package internal

import (
	"Stant/LestaGamesInternship/internal/views"
	"bufio"
	"log"
	"maps"
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

		uniqueWords := map[string]uint64{}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Text()
			uniqueWords[word] = uniqueWords[word] + 1
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			return
		}

		for word, amount := range maps.All(uniqueWords) {
			log.Printf("%q: %d\n", word, amount)
		}
	})
}
