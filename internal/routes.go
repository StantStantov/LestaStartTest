package internal

import (
	"Stant/LestaGamesInternship/internal/models"
	"Stant/LestaGamesInternship/internal/services"
	"Stant/LestaGamesInternship/internal/stores"
	"Stant/LestaGamesInternship/internal/views"
	"context"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
)

func HandleIndexGet(termStore stores.TermStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})
}

func HandleIndexPost(termStore stores.TermStore) http.HandlerFunc {
	MaxTableLength := 50
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		text, err := services.ProcessReaderToTerms(file)
		if err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		totalAmount := uint64(len(text))
		uniqueWords := services.GetTermFrequency(text)

		length, err := termStore.CountAll()
		if err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			http.Error(w, "Failed to access database", http.StatusInternalServerError)
			return
		}
		for id := range length {
			termStore.Delete(length - id - 1)
		}
		for word, amount := range maps.All(uniqueWords) {
			term := models.NewTerm(word, amount, services.CalculateIdf(totalAmount, amount))
			if err := termStore.Create(term); err != nil {
				log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
				http.Error(w, "Failed to access database", http.StatusInternalServerError)
				return
			}
		}

		table, err := termStore.ReadAll()
		if err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			http.Error(w, "Failed to access database", http.StatusInternalServerError)
			return
		}
		slices.SortFunc(table, compareRowsByIdf)
		if len(table) > MaxTableLength {
			table = table[:MaxTableLength]
		}

		renderTable(table, w, r.Context())
	})
}

func compareRowsByIdf(a models.Term, b models.Term) int {
	if a.Idf() < b.Idf() {
		return 1
	} else if a.Idf() > b.Idf() {
		return -1
	}
	return 0
}

func renderTable(table []models.Term, w http.ResponseWriter, rCtx context.Context) {
	tableViewModel := make([]views.TableRowViewModel, len(table))
	for i, row := range table {
		rowViewModel := views.TableRowViewModel{
			Word: row.Word(),
			Tf:   strconv.FormatUint(row.Frequency(), 10),
			Idf:  strconv.FormatFloat(row.Idf(), 'G', -1, 64),
		}
		tableViewModel[i] = rowViewModel
	}

	w.WriteHeader(http.StatusOK)
	views.Table(tableViewModel).Render(rCtx, w)
}
