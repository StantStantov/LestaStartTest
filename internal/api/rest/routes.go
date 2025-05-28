package rest

import (
	"Stant/LestaGamesInternship/internal/app/services/metricService"
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func SetupRestRouter(router *http.ServeMux, metricsStore stores.MetricStore) {
	router.Handle("GET /api/status", HandleStatusGet())
	router.Handle("GET /api/metrics", HandleGetMetrics(metricsStore))
	router.Handle("GET /api/version", HandleGetVerstion())
}

func HandleStatusGet() http.HandlerFunc {
	const status = `{"status": "OK"}`

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, status)
	})
}

func HandleGetMetrics(metricsStore stores.MetricStore) http.HandlerFunc {
	metricsJson := metricService.AppMetrics{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filesMetrics, err := metricsStore.ReadAllByName(models.FilesProcessed)
		if err != nil {
			log.Printf("Internal/rest.HandleGetMetrics: [%v]", err)
			http.Error(w, "Failed to access database", http.StatusInternalServerError)
			return
		}
		timeMetrics, err := metricsStore.ReadAllByName(models.TimeProcessed)
		if err != nil {
			log.Printf("Internal/rest.HandleGetMetrics: [%v]", err)
			http.Error(w, "Failed to access database", http.StatusInternalServerError)
			return
		}

		if len(filesMetrics) != 0 {
			filesProcessedCount, err := metricService.SumValues(filesMetrics)
			if err != nil {
				log.Printf("Internal/rest.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}
			latestFileProcessed, err := metricService.FindMaxByTimestamp(filesMetrics)
			if err != nil {
				log.Printf("Internal/rest.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}

			metricsJson.FilesProcessed = uint64(filesProcessedCount)
			metricsJson.LatestFileProcessed = &latestFileProcessed
		}

		if len(timeMetrics) != 0 {
			timeProcessedCount, err := metricService.SumValues(timeMetrics)
			if err != nil {
				log.Printf("Internal/rest.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}
			minTimeProcessed, err := metricService.FindMinByValue(timeMetrics)
			if err != nil {
				log.Printf("Internal/rest.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}
			maxTimeProcessed, err := metricService.FindMaxByValue(timeMetrics)
			if err != nil {
				log.Printf("Internal/rest.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}

			metricsJson.MinTimeProcessed = minTimeProcessed
			metricsJson.AvgTimeProcessed = (maxTimeProcessed - minTimeProcessed) / timeProcessedCount
			metricsJson.MaxTimeProcessed = maxTimeProcessed
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(metricsJson)
	})
}

func HandleGetVerstion() http.HandlerFunc {
	info, _ := debug.ReadBuildInfo()
	version := fmt.Sprintf(`{"version": "%s"}`, info.Main.Version)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, version)
	})
}
