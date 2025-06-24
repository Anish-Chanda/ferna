package model

import "time"

type User struct {
	ID        int64
	Email     string
	PassHash  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Species struct {
	ID                       int64     `json:"id"`
	CommonName               string    `json:"common_name"`
	ScientificName           string    `json:"scientific_name,omitempty"`
	DefaultWateringFrequency int       `json:"default_watering_frequency_days"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}
