package model

type CleanerStat struct {
    CleanerName   string  `json:"cleaner_name"`
    TotalLogs     int     `json:"total_logs"`
    TotalMinutes  float64 `json:"total_minutes"` // Durasi dalam menit
}