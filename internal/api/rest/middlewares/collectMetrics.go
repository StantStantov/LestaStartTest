package middlewares

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"net/http"
	"time"
)

func NewCollectMetricsMiddleware(metricStore stores.MetricStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTimestamp := time.Now()

			wrappedWriter := wrapResponseWriter(w)
			next.ServeHTTP(wrappedWriter, r)

			endTimestamp := time.Now()
			duration := endTimestamp.Sub(startTimestamp)

			if wrappedWriter.status == http.StatusOK {
				metricStore.Track(r.Context(), models.NewMetric(endTimestamp, models.FilesProcessed, 1))
				metricStore.Track(r.Context(), models.NewMetric(endTimestamp, models.TimeProcessed, duration.Seconds()))
			}
		})
	}
}
