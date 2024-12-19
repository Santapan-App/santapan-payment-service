package domain

import (
	"time"
)

// Payment represents a payment entity
type Payment struct {
	ID          int64     `json:"id"`
	ReferenceID string    `json:"reference_id"`
	SessionID   string    `json:"session_id"`
	UserID      int64     `json:"user_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaymentBody struct {
	Amount float64   `json:"amount" validate:"required"`
	Name   []string  `json:"name" validate:"required,dive,required"` // Ensures each name is non-empty
	Qty    []int64   `json:"qty" validate:"required,dive,gt=0"`      // Ensures each qty is > 0
	Price  []float64 `json:"price" validate:"required,dive,gt=0"`    // Ensures each price is > 0
}

type IPaymuResponse struct {
	Status  int         `json:"Status"`  // The status code of the response, typically used to indicate success or failure.
	Succees bool        `json:"Succees"` // A boolean value indicating whether the request was successful.
	Data    *IPaymuData `json:"Data"`    // IPaymuData is now a pointer, making it nullable.
	Message string      `json:"Message"` // A message providing more details about the response.
}

type IPaymuData struct {
	SessionID string `json:"SessionID"`
	Url       string `json:"Url"`
}

type IPaymuCallback struct {
	TrxID       string `json:"trx_id"`
	ReferenceID string `json:"reference_id"`
	SessionID   string `json:"sid"`
	Status      string `json:"status"`
	StatusCode  string `json:"status_code"`
}
