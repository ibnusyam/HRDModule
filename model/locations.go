package model

// --- Location Type ---
type LocationTypeIt struct {
	LocationTypeID int    `json:"location_type_id"`
	TypeName       string `json:"type_name"`
	Description    string `json:"description"`
	SiteID         int    `json:"site_id"`
}

// --- Location ---
type LocationIt struct {
	LocationID     int    `json:"location_id"`
	LocationName   string `json:"location_name"`
	LocationTypeID int    `json:"location_type_id"`
	TypeName       string `json:"type_name,omitempty"` // Field tambahan hasil JOIN (untuk display)
	SiteID         int    `json:"site_id"`
}