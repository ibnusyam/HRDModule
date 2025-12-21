package repository

import (
	"HRD/model" // Sesuaikan nama module
	"database/sql"
)

type LocationRepository struct {
	DB *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{DB: db}
}

// ================= LOCATION TYPES =================

func (r *LocationRepository) GetLocationTypes() ([]model.LocationTypeIt, error) {
	query := `SELECT location_type_id, type_name, COALESCE(description, ''), site_id FROM location_types ORDER BY location_type_id ASC`
	rows, err := r.DB.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var result []model.LocationTypeIt
	for rows.Next() {
		var d model.LocationTypeIt
		if err := rows.Scan(&d.LocationTypeID, &d.TypeName, &d.Description, &d.SiteID); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func (r *LocationRepository) CreateLocationType(input model.LocationTypeIt) error {
	query := `INSERT INTO location_types (type_name, description, site_id) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, input.TypeName, input.Description, input.SiteID)
	return err
}

func (r *LocationRepository) UpdateLocationType(id int, input model.LocationTypeIt) error {
	query := `UPDATE location_types SET type_name=$1, description=$2, site_id=$3 WHERE location_type_id=$4`
	_, err := r.DB.Exec(query, input.TypeName, input.Description, input.SiteID, id)
	return err
}

func (r *LocationRepository) DeleteLocationType(id int) error {
	_, err := r.DB.Exec("DELETE FROM location_types WHERE location_type_id=$1", id)
	return err
}

// ================= LOCATIONS =================

func (r *LocationRepository) GetLocations() ([]model.LocationIt, error) {
	// JOIN supaya dapat TypeName
	query := `
		SELECT l.location_id, l.location_name, l.location_type_id, t.type_name, l.site_id
		FROM locations l
		LEFT JOIN location_types t ON l.location_type_id = t.location_type_id
		ORDER BY l.location_id ASC
	`
	rows, err := r.DB.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var result []model.LocationIt
	for rows.Next() {
		var d model.LocationIt
		if err := rows.Scan(&d.LocationID, &d.LocationName, &d.LocationTypeID, &d.TypeName, &d.SiteID); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func (r *LocationRepository) CreateLocation(input model.Location) error {
	query := `INSERT INTO locations (location_name, location_type_id, site_id) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, input.LocationName, input.LocationTypeID, input.SiteID)
	return err
}

func (r *LocationRepository) UpdateLocation(id int, input model.Location) error {
	query := `UPDATE locations SET location_name=$1, location_type_id=$2, site_id=$3 WHERE location_id=$4`
	_, err := r.DB.Exec(query, input.LocationName, input.LocationTypeID, input.SiteID, id)
	return err
}

func (r *LocationRepository) DeleteLocation(id int) error {
	_, err := r.DB.Exec("DELETE FROM locations WHERE location_id=$1", id)
	return err
}