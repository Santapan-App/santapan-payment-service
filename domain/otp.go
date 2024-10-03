package domain

import (
	"time"
)

type Otp struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Retry     int16     `json:"retry"`
	Type      string    `json:"type"`
	UserId    int64     `json:"user_id"`
	DeviceId  string    `json:"device_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
