package web

import (
	"Stant/LestaGamesInternship/internal/api/web/middlewares"
	"Stant/LestaGamesInternship/internal/api/web/views"
	"Stant/LestaGamesInternship/internal/app/services"
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
)

func SetupWebRouter(router *http.ServeMux, metricsStore stores.MetricStore, termStore stores.TermStore) {
	collectMetrics := middlewares.NewCollectMetricsMiddleware(metricsStore)

	router.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	router.Handle("GET /", HandleGetIndex(termStore))
	router.Handle("POST /", collectMetrics(HandlePostIndex(termStore)))
}

func HandleGetIndex(termStore stores.TermStore) http.HandlerFunc {
	MaxTableLength := 50
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		tableViewModel := renderTable(table)

		w.WriteHeader(http.StatusOK)
		views.Index(tableViewModel).Render(r.Context(), w)
	})
}

func HandlePostIndex(termStore stores.TermStore) http.HandlerFunc {
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
			if err := termStore.Delete(length - id - 1); err != nil {
				log.Printf("Internal/routes.HandleIndexPost: [%v]", err)
				http.Error(w, "Failed to access database", http.StatusInternalServerError)
				return
			}
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

		tableViewModel := renderTable(table)
		w.WriteHeader(http.StatusOK)
		views.Table(tableViewModel).Render(r.Context(), w)
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

func renderTable(table []models.Term) []views.TableRowViewModel {
	tableViewModel := make([]views.TableRowViewModel, len(table))
	for i, row := range table {
		rowViewModel := views.TableRowViewModel{
			Word: row.Word(),
			Tf:   strconv.FormatUint(row.Frequency(), 10),
			Idf:  strconv.FormatFloat(row.Idf(), 'G', -1, 64),
		}
		tableViewModel[i] = rowViewModel
	}
	return tableViewModel
}
