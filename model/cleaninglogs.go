package model

import (
	"database/sql"
	"time"
)

type CleaningLog struct {
    LogID            int64          `json:"log_id"`
    CleanerName      string         `json:"cleaner_name"`
    LocationID       int            `json:"location_id"`  
    LocationTypeID   int            `json:"location_type_id"` 
    LocationName     string         `json:"location_name,omitempty"`
    LocationTypeName string         `json:"location_type_name,omitempty"`
    StartTime        time.Time      `json:"start_time"`
    EndTime          sql.NullTime   `json:"end_time"`
    ImageBeforeURL   string         `json:"image_before_url"`
    ImageAfterURL    sql.NullString `json:"image_after_url"`
    Notes            sql.NullString `json:"notes"`
    SiteID           int            `json:"site_id"`
}

type CreateFullLogInput struct {
	CleanerName      string `json:"cleaner_name"`
	LocationID       int    `json:"location_id"`
	LocationTypeID   int    `json:"location_type_id"`
	StartTimeStr     string `json:"start_time"`
	EndTimeStr       string `json:"end_time"`
	Notes            string `json:"notes"`
	SiteID           int    `json:"site_id"`
}

type LocationType struct {
	LocationTypeID int    `json:"location_type_id"`
	TypeName       string `json:"type_name"`
}

type Location struct {
	LocationID     int    `json:"location_id"`
	LocationName   string `json:"location_name"`
	LocationTypeID int    `json:"location_type_id"`
}

type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	Limit        int `json:"limit"`
}

type PaginatedResponse struct {
	Data []CleaningLog  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}