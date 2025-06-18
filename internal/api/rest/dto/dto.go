package dto

import "time"

type Collection struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Documents []Document `json:"documents"`
}

type Document struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DocumentWithData struct {
	Document
	Data string `json:"data"`
}

type Term struct {
	Word string  `json:"word"`
	Tf   uint64  `json:"tf"`
	Idf  float64 `json:"idf"`
}

type AppMetrics struct {
	FilesProcessed      uint64     `json:"files_processed"`
	MinTimeProcessed    float64    `json:"min_time_processed"`
	AvgTimeProcessed    float64    `json:"avg_time_processed"`
	MaxTimeProcessed    float64    `json:"max_time_processed"`
	LatestFileProcessed *time.Time `json:"latest_file_processed_timestamp"`
}
