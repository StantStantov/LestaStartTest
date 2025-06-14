package models

import "time"

type MetricName string

const (
	FilesProcessed MetricName = "files_processed"
	TimeProcessed  MetricName = "time_processed"
)

type Metric struct {
	timestamp time.Time
	name      MetricName
	value     float64
}

func NewMetric(timestamp time.Time, metricName MetricName, value float64) Metric {
	return Metric{
		timestamp: timestamp,
		name:      metricName,
		value:     value,
	}
}

func (m Metric) Name() MetricName {
	return m.name
}

func (m Metric) Timestamp() time.Time {
	return m.timestamp
}

func (m Metric) Value() float64 {
	return m.value
}
