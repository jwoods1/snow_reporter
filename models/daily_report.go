package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// DailyReport is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type DailyReport struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Date       time.Time `json:"date" db:"date"`
	StationID  uuid.UUID `json:"station_id" db:"station_id"`
	SnowDepth  float64   `json:"snow_depth" db:"snow_depth"`
	SnowChange float64   `json:"snow_change" db:"snow_change"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Station    *Station  `json:"station,omitempty" belongs_to:"station"`
}

// String is not required by pop and may be deleted
func (d DailyReport) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// DailyReports is not required by pop and may be deleted
type DailyReports []DailyReport

// String is not required by pop and may be deleted
func (d DailyReports) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *DailyReport) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *DailyReport) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *DailyReport) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
