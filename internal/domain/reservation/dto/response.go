package dto

import "time"

type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type MyReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    time.Time   `json:"created_at"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
