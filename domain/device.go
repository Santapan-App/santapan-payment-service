package domain

import "time"

// Device represents a device entity in the system
type Device struct {
	ID        int64      `json:"id"`                   // Maps to id BIGINT
	Name      string     `json:"name"`                 // Maps to name VARCHAR(255)
	Brand     string     `json:"brand"`                // Maps to brand VARCHAR(255)
	UniqueID  string     `json:"unique_id"`            // Maps to unique_id VARCHAR(255)
	UserID    int64      `json:"user_id"`              // Maps to user_id BIGINT NOT NULL
	IPAddress string     `json:"ip_address"`           // Maps to ip_address VARCHAR(255)
	CreatedAt time.Time  `json:"created_at"`           // Maps to created_at TIMESTAMP WITH TIME ZONE
	UpdatedAt time.Time  `json:"updated_at"`           // Maps to updated_at TIMESTAMP WITH TIME ZONE
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Maps to deleted_at TIMESTAMP WITH TIME ZONE, nullable
}

// DeviceHeaderInformation represents the header information of a device
type DeviceHeaderInformation struct {
	DeviceID    string
	DeviceName  string
	DeviceModel string
	DeviceBrand string
	IPAddress   string
}
