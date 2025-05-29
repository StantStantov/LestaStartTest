package metricService

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"cmp"
	"fmt"
	"slices"
	"time"
)

func SumValues(metrics []models.Metric) float64 {
	var sumValues float64 = 0
	for _, metric := range metrics {
		sumValues += metric.Value()
	}

	return sumValues
}

func FindMaxByValue(metrics []models.Metric) (float64, error) {
	if len(metrics) == 0 {
		return 0, fmt.Errorf("Services/metricService.FindMaxByValue: [metrics are empty]")
	}

	return slices.MaxFunc(metrics, compareByValue).Value(), nil
}

func FindMaxByTimestamp(metrics []models.Metric) (time.Time, error) {
	if len(metrics) == 0 {
		return time.Time{}, fmt.Errorf("Services/metricService.FindMaxByTimestamp: [metrics are empty]")
	}

	return slices.MaxFunc(metrics, compareByTimestamp).Timestamp(), nil
}

func FindMinByValue(metrics []models.Metric) (float64, error) {
	if len(metrics) == 0 {
		return 0, fmt.Errorf("Services/metricService.FindMinByValue: [metrics are empty]")
	}

	return slices.MinFunc(metrics, compareByValue).Value(), nil
}

func FindMinByTimestamp(metrics []models.Metric) (time.Time, error) {
	if len(metrics) == 0 {
		return time.Time{}, fmt.Errorf("Services/metricService.FindMinByTimestamp: [metrics are empty]")
	}

	return slices.MinFunc(metrics, compareByTimestamp).Timestamp(), nil
}

func compareByValue(a, b models.Metric) int {
	return cmp.Compare(a.Value(), b.Value())
}

func compareByTimestamp(a, b models.Metric) int {
	return a.Timestamp().Compare(b.Timestamp())
}
