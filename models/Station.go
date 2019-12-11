package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Station is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Station struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	StationID    int           `json:"station_id" db:"station_id"`
	Name         string        `json:"name" db:"name"`
	State        string        `json:"state" db:"state"`
	Elevation    int           `json:"elevation" db:"elevation"`
	Lat          float64       `json:"lat" db:"lat"`
	Long         float64       `json:"long" db:"long"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
	DailyReports []DailyReport `json:"daily_reports,omitempty" has_many:"daily_reports"`
}

// String is not required by pop and may be deleted
func (s Station) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Stations is not required by pop and may be deleted
type Stations []Station

// String is not required by pop and may be deleted
func (s Stations) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Station) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Station) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Station) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
