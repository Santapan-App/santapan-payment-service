package domain

import "time"

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        int64     `json:"id"`         // Unique identifier for the cart item
	Name      string    `json:"name"`       // Name of the item
	Quantity  int       `json:"quantity"`   // Quantity of the item in the cart
	Price     float64   `json:"price"`      // Price per unit of the item
	Subtotal  float64   `json:"subtotal"`   // Calculated as Quantity * Price
	CreatedAt time.Time `json:"created_at"` // Timestamp of when the item was added to the cart
	UpdatedAt time.Time `json:"updated_at"` // Timestamp of the last update to the item
}
