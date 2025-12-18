package repository

import (
	"HRD/model"
	"database/sql"
)

type DashboardRepository struct {
    DB *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
    return &DashboardRepository{DB: db}
}

func (repo *DashboardRepository) GetCleanerStatsByMonth(siteID, month, year int) ([]model.CleanerStat, error) {
    // PERBAIKAN: Menggunakan start_time sebagai acuan waktu filter
    // Durasi dihitung hanya jika end_time tidak NULL.
    query := `
        SELECT 
            cleaner_name,
            COUNT(log_id) as total_logs,
            COALESCE(SUM(EXTRACT(EPOCH FROM (end_time - start_time))/60), 0) as total_minutes
        FROM cleaning_logs
        WHERE site_id = $1 
          AND EXTRACT(MONTH FROM start_time) = $2
          AND EXTRACT(YEAR FROM start_time) = $3
        GROUP BY cleaner_name
        ORDER BY total_logs DESC
    `

    rows, err := repo.DB.Query(query, siteID, month, year)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []model.CleanerStat
    for rows.Next() {
        var s model.CleanerStat
        if err := rows.Scan(&s.CleanerName, &s.TotalLogs, &s.TotalMinutes); err != nil {
            return nil, err
        }
        stats = append(stats, s)
    }

    return stats, nil
}