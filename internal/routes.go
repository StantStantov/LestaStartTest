package internal

import (
	"Stant/LestaGamesInternship/internal/services"
	"Stant/LestaGamesInternship/internal/views"
	"context"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
)

func HandleIndexGet() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})
}

func HandleIndexPost() http.HandlerFunc {
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

		table := make([]tableRow, 0, totalAmount)
		for word, amount := range maps.All(uniqueWords) {
			table = append(table, tableRow{word, amount, services.CalculateIdf(totalAmount, amount)})
		}
		slices.SortFunc(table, compareRowsByIdf)
		if len(table) > MaxTableLength {
			table = table[:MaxTableLength]
		}

		renderTable(table, w, r.Context())
	})
}

type tableRow struct {
	word string
	tf   uint64
	idf  float64
}

func compareRowsByIdf(a tableRow, b tableRow) int {
	if a.idf < b.idf {
		return 1
	} else if a.idf > b.idf {
		return -1
	}
	return 0
}

func renderTable(table []tableRow, w http.ResponseWriter, rCtx context.Context) {
	tableViewModel := make([]views.TableRowViewModel, len(table))
	for i, row := range table {
		rowViewModel := views.TableRowViewModel{
			Word: row.word,
			Tf:   strconv.FormatUint(row.tf, 10),
			Idf:  strconv.FormatFloat(row.idf, 'G', -1, 64),
		}
		tableViewModel[i] = rowViewModel
	}

	w.WriteHeader(http.StatusOK)
	views.Table(tableViewModel).Render(rCtx, w)
}
