package handlers

import (
	"Stant/LestaGamesInternship/internal/api/rest/dto"
	"Stant/LestaGamesInternship/internal/app/config"
	"Stant/LestaGamesInternship/internal/app/services/mtrcserv"
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// @Summary Получить статус приложения
// @Description Получает текущий статус приложения.
// @Tags Общее
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/status [get]
func HandleGetStatus() http.HandlerFunc {
	const status = `{"status": "OK"}`

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, status)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// @Summary Получить метрики приложения
// @Description Получает текущие метрики приложения.
// @Tags Общее
// @Produce json
// @Success 200 {object} dto.AppMetrics
// @Router /api/metrics [get]
func HandleGetMetrics(metricsStore stores.MetricStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metricsJson := dto.AppMetrics{}

		filesMetrics, err := metricsStore.FindAllByName(r.Context(), models.FilesProcessed)
		if err != nil && errors.Unwrap(errors.Unwrap(err)) != pgx.ErrNoRows {
			log.Printf("handlers/info.HandleGetMetrics: [%v]", err)
			http.Error(w, "Failed to access database", http.StatusInternalServerError)
			return
		}
		timeMetrics, err := metricsStore.FindAllByName(r.Context(), models.TimeProcessed)
		if err != nil && errors.Unwrap(errors.Unwrap(err)) != pgx.ErrNoRows {
			log.Printf("handlers/info.HandleGetMetrics: [%v]", err)
			http.Error(w, "Failed to access database", http.StatusInternalServerError)
			return
		}

		if len(filesMetrics) != 0 {
			filesProcessedCount := mtrcserv.SumValues(filesMetrics)
			latestFileProcessed, err := mtrcserv.FindMaxByTimestamp(filesMetrics)
			if err != nil {
				log.Printf("handlers/info.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}

			metricsJson.FilesProcessed = uint64(filesProcessedCount)
			metricsJson.LatestFileProcessed = &latestFileProcessed
		}

		if len(timeMetrics) != 0 {
			timeProcessedCount := mtrcserv.SumValues(timeMetrics)
			minTimeProcessed, err := mtrcserv.FindMinByValue(timeMetrics)
			if err != nil {
				log.Printf("handlers/info.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}
			maxTimeProcessed, err := mtrcserv.FindMaxByValue(timeMetrics)
			if err != nil {
				log.Printf("handlers/info.HandleGetMetrics: [%v]", err)
				http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
				return
			}

			metricsJson.MinTimeProcessed = minTimeProcessed
			metricsJson.AvgTimeProcessed = (maxTimeProcessed - minTimeProcessed) / timeProcessedCount
			metricsJson.MaxTimeProcessed = maxTimeProcessed
		}

		if err := json.NewEncoder(w).Encode(metricsJson); err != nil {
			log.Printf("handlers/info.HandleGetMetrics: [%v]", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// @Summary Получить версию приложения
// @Description Получает текущую версию приложения.
// @Tags Общее
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/version [get]
func HandleGetVersion(config *config.AppConfig) http.HandlerFunc {
	version := fmt.Sprintf(`{"version": "%s"}`, config.Version())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, version)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
