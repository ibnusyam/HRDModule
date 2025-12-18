package repository

import (
	"HRD/model"
	"database/sql"
	"fmt"
)

type CleaningLogsRepository struct {
	DB *sql.DB
}

func NewCleaningLogsRepository(db *sql.DB) *CleaningLogsRepository {
	return &CleaningLogsRepository{DB: db}
}

func buildFilterQuery(siteID, locationID, typeID int, cleanerName, dateStr string) (string, []interface{}) {
	query := `
        FROM cleaning_logs cl
        JOIN locations l ON cl.location_id = l.location_id
        JOIN location_types lt ON cl.location_type_id = lt.location_type_id 
        WHERE cl.site_id = $1
    `
	args := []interface{}{siteID}
	argCounter := 2

	if locationID > 0 {
		query += fmt.Sprintf(" AND cl.location_id = $%d", argCounter)
		args = append(args, locationID)
		argCounter++
	}

	if typeID > 0 {
		query += fmt.Sprintf(" AND cl.location_type_id = $%d", argCounter)
		args = append(args, typeID)
		argCounter++
	}

	if cleanerName != "" {
		query += fmt.Sprintf(" AND cl.cleaner_name ILIKE $%d", argCounter)
		args = append(args, "%"+cleanerName+"%")
		argCounter++
	}

	if dateStr != "" {
		query += fmt.Sprintf(" AND TO_CHAR(cl.start_time, 'YYYY-MM') = $%d", argCounter)
		args = append(args, dateStr)
		argCounter++
	}

	return query, args
}

func (repo *CleaningLogsRepository) CountLogs(siteID, locationID, typeID int, cleanerName, dateStr string) (int, error) {
	baseQuery, args := buildFilterQuery(siteID, locationID, typeID, cleanerName, dateStr)
	query := "SELECT COUNT(*) " + baseQuery

	var count int
	err := repo.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *CleaningLogsRepository) GetAllCleaningsLogs(siteID, locationID, typeID, limit, offset int, cleanerName, dateStr string) ([]model.CleaningLog, error) {
    baseQuery, args := buildFilterQuery(siteID, locationID, typeID, cleanerName, dateStr)

    query := `
        SELECT 
            cl.log_id, cl.cleaner_name, l.location_name, lt.type_name,
            cl.start_time, cl.end_time, cl.image_before_url, cl.image_after_url, cl.notes,
            cl.location_id, cl.location_type_id
    ` + baseQuery + " ORDER BY cl.start_time DESC LIMIT $" + fmt.Sprint(len(args)+1) + " OFFSET $" + fmt.Sprint(len(args)+2)

    args = append(args, limit, offset)
    rows, err := repo.DB.Query(query, args...)
    if err != nil { return nil, err }
    defer rows.Close()

    var logs []model.CleaningLog
    for rows.Next() {
        log := model.CleaningLog{}
        err := rows.Scan(
            &log.LogID, &log.CleanerName, &log.LocationName, &log.LocationTypeName,
            &log.StartTime, &log.EndTime, &log.ImageBeforeURL, &log.ImageAfterURL, &log.Notes,
            &log.LocationID, &log.LocationTypeID, // Tambahkan Scan untuk field ID
        )
        if err != nil { return nil, err }
        logs = append(logs, log)
    }
    return logs, nil
}

func (repo *CleaningLogsRepository) CreateFullLog(log *model.CleaningLog) error {
	query := `
        INSERT INTO cleaning_logs (
            cleaner_name, location_id, location_type_id, 
            start_time, end_time, 
            image_before_url, image_after_url, notes, site_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING log_id
    `

	err := repo.DB.QueryRow(
		query,
		log.CleanerName,
		log.LocationID,
		log.LocationTypeID,
		log.StartTime,
		log.EndTime,
		log.ImageBeforeURL,
		log.ImageAfterURL,
		log.Notes,
		log.SiteID,
	).Scan(&log.LogID)

	if err != nil {
		return fmt.Errorf("failed to insert full cleaning log: %w", err)
	}

	return nil
}

func (repo *CleaningLogsRepository) GetAllLocationTypes() ([]model.LocationType, error) {
	query := `SELECT location_type_id, type_name FROM location_types ORDER BY type_name ASC`
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []model.LocationType
	for rows.Next() {
		var t model.LocationType
		if err := rows.Scan(&t.LocationTypeID, &t.TypeName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func (repo *CleaningLogsRepository) GetAllLocations() ([]model.Location, error) {
	query := `SELECT location_id, location_name, location_type_id FROM locations WHERE is_active = true ORDER BY location_name ASC`
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locs []model.Location
	for rows.Next() {
		var l model.Location
		if err := rows.Scan(&l.LocationID, &l.LocationName, &l.LocationTypeID); err != nil {
			return nil, err
		}
		locs = append(locs, l)
	}
	return locs, nil
}