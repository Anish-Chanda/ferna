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

// Plant represents a user’s individual plant.
type Plant struct {
	ID                    int64      `json:"id"`
	UserID                int64      `json:"user_id"`
	SpeciesID             int64      `json:"species_id"`
	Nickname              *string    `json:"nickname"`
	ImageURL              *string    `json:"image_url"`
	WateringFrequencyDays int        `json:"watering_frequency_days"`
	LastWateredAt         *time.Time `json:"last_watered_at"`
	Note                  *string    `json:"note"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
