package metricService

import "time"

type AppMetrics struct {
	FilesProcessed      uint64     `json:"files_processed"`
	MinTimeProcessed    float64    `json:"min_time_processed"`
	AvgTimeProcessed    float64    `json:"avg_time_processed"`
	MaxTimeProcessed    float64    `json:"max_time_processed"`
	LatestFileProcessed *time.Time `json:"latest_file_processed_timestamp"`
}
