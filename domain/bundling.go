package domain

import (
	"time"
)

// Menu represents the menu table
type Menu struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"` // Allows NULL values
	Price       int64     `json:"price"`
	ImageURL    string    `json:"image_url"` // Allows NULL values
	Nutrition   []byte    `json:"nutrition"` // JSON data stored as []byte
	Features    []byte    `json:"features"`  // JSON data stored as []byte
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Bundling represents the bundling table
type Bundling struct {
	ID           int64     `json:"id"`
	BundlingType string    `json:"bundling_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BundlingMenu represents the bundling_menu table
type BundlingMenu struct {
	ID              int64     `json:"id"`
	BundlingID      int64     `json:"bundling_id"`
	MenuID          int64     `json:"menu_id"`
	DayNumber       int       `json:"day_number"`
	MealDescription string    `json:"meal_description"` // Allows NULL values
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
