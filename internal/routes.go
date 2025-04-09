package internal

import (
	"Stant/LestaGamesInternship/internal/views"
	"bufio"
	"context"
	"log"
	"maps"
	"math"
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

		var totalAmount uint64 = 0
		uniqueWords := map[string]uint64{}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Text()
			uniqueWords[word] = uniqueWords[word] + 1
			totalAmount++
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		table := make([]tableRow, 0, totalAmount)
		for word, amount := range maps.All(uniqueWords) {
			table = append(table, tableRow{word, amount, calculateIdf(totalAmount, amount)})
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

func calculateIdf(wordsAmount uint64, wordAmount uint64) float64 {
	return math.Log10(float64(wordsAmount) / float64(wordAmount))
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
