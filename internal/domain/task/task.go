package task

import (
	"encoding/json"
	"time"
)

type RepeatType string

const (
	RepeatDaily    RepeatType = "daily"
	RepeatMonthly  RepeatType = "monthly"
	RepeatSpecific RepeatType = "specific"
	RepeatParity   RepeatType = "parity"
)

type RepeatRule struct {
	Type         RepeatType `json:"type"`
	Interval     int        `json:"interval,omitempty"`
	DayOfMonth   int        `json:"day_of_month,omitempty"`
	SpecificDate string     `json:"specific_date,omitempty"`
	Parity       string     `json:"parity,omitempty"`
}

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID            int64       `json:"id"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Status        Status      `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	StartDateTime *time.Time  `json:"start_datetime,omitempty"`
	EndDateTime   *time.Time  `json:"end_datetime,omitempty"`
	RepeatRule    *RepeatRule `json:"repeat_rule"`
}

func (r *RepeatRule) Valid() bool {
	if r == nil {
		return true
	}

	switch r.Type {
	case RepeatDaily:
		return r.Interval > 0
	case RepeatMonthly:
		return r.DayOfMonth >= 1 && r.DayOfMonth <= 30
	case RepeatSpecific:
		_, err := time.Parse("2006-01-02", r.SpecificDate)
		return err == nil
	case RepeatParity:
		return r.Parity == "even" || r.Parity == "odd"
	default:
		return false
	}
}

func (r *RepeatRule) MarshalJSON() ([]byte, error) {
	type Alias RepeatRule
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *RepeatRule) UnmarshalJSON(data []byte) error {
	type Alias RepeatRule
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	return json.Unmarshal(data, aux)
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
