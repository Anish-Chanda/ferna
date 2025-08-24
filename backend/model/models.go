package model

import "time"

// AuthProvider represents the authentication provider
type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGoogle AuthProvider = "google"
)

// LightPreference represents light requirements for plants
type LightPreference string

const (
	LightPreferenceLow            LightPreference = "low"
	LightPreferenceMedium         LightPreference = "medium"
	LightPreferenceBrightIndirect LightPreference = "bright_indirect"
	LightPreferenceBrightDirect   LightPreference = "bright_direct"
)

// Toxicity represents toxicity level for pets
type Toxicity string

const (
	ToxicityToxicToPets Toxicity = "toxic_to_pets"
	ToxicityNonToxic    Toxicity = "non_toxic"
)

// TaskType represents the type of plant care task
type TaskType string

const (
	TaskTypeWatering   TaskType = "watering"
	TaskTypeFertilizer TaskType = "fertilizer"
)

// EventType represents the type of care event
type EventType string

const (
	EventTypeWatering   EventType = "watering"
	EventTypeFertilizer EventType = "fertilizer"
	EventTypeRepotting  EventType = "repotting"
	EventTypePruning    EventType = "pruning"
	EventTypeOther      EventType = "other"
)

type User struct {
	ID             int64        `json:"id"`
	Email          string       `json:"email"`
	PasswordHash   *string      `json:"password_hash"`
	AuthProvider   AuthProvider `json:"auth_provider"`
	ProviderUserID *string      `json:"provider_user_id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type Species struct {
	ID                            int64           `json:"id"`
	CommonName                    string          `json:"common_name"`
	ScientificName                *string         `json:"scientific_name"`
	LightPreference               LightPreference `json:"light_pref"`
	DefaultWaterIntervalDays      int             `json:"default_water_interval_days"`
	DefaultFertilizerIntervalDays int             `json:"default_fertilizer_interval_days"`
	Toxicity                      Toxicity        `json:"toxicity"`
	CareNotes                     *string         `json:"care_notes"`
	CareNotesSource               *string         `json:"care_notes_source"`
	CreatedAt                     time.Time       `json:"created_at"`
	UpdatedAt                     time.Time       `json:"updated_at"`
}

type Location struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserPlant struct {
	ID                             int64     `json:"id"`
	UserID                         int64     `json:"user_id"`
	SpeciesID                      int64     `json:"species_id"`
	Nickname                       *string   `json:"nickname"`
	ImageURL                       *string   `json:"image_url"`
	Notes                          *string   `json:"notes"`
	WaterIntervalDaysOverride      *int      `json:"water_interval_days_override"`
	FertilizerIntervalDaysOverride *int      `json:"fertilizer_interval_days_override"`
	LocationID                     *int64    `json:"location_id"`
	CreatedAt                      time.Time `json:"created_at"`
	UpdatedAt                      time.Time `json:"updated_at"`
}

type PlantTask struct {
	ID            int64      `json:"id"`
	PlantID       int64      `json:"plant_id"`
	TaskType      TaskType   `json:"task_type"`
	SnoozedUntil  *time.Time `json:"snoozed_until"`
	IntervalDays  int        `json:"interval_days"`
	ToleranceDays int        `json:"tolerance_days"`
	NextDueAt     *time.Time `json:"next_due_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CareEvent struct {
	ID         int64     `json:"id"`
	PlantID    int64     `json:"plant_id"`
	TaskID     *int64    `json:"task_id"`
	EventType  EventType `json:"event_type"`
	HappenedAt time.Time `json:"happened_at"`
	Notes      *string   `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
