package model

import "fmt"

type RecordStatus string

const (
	StatusDraft    RecordStatus = "draft"
	StatusReviewed RecordStatus = "reviewed"
	StatusApproved RecordStatus = "approved"
	StatusArchived RecordStatus = "archived"
)

type Record struct {
	ID         string       `json:"id"`
	Classroom  string       `json:"classroom"`
	Student    string       `json:"student"`
	Gesture    string       `json:"gesture"`
	Particle   string       `json:"particle"`
	Intensity  int          `json:"intensity"`
	Sequence   int          `json:"sequence"`
	Status     RecordStatus `json:"status"`
	CreatedBy  string       `json:"created_by"`
	Reviewer   string       `json:"reviewer,omitempty"`
	ReviewedAt string       `json:"reviewed_at,omitempty"`
	Version    int          `json:"version"`
	Notes      string       `json:"notes,omitempty"`
}

func (r Record) Validate() error {
	if r.ID == "" || r.Classroom == "" || r.Student == "" {
		return fmt.Errorf("record identity is required")
	}
	if r.Gesture == "" || r.Particle == "" {
		return fmt.Errorf("gesture and particle are required")
	}
	if r.Intensity < 1 || r.Intensity > 10 {
		return fmt.Errorf("intensity must be between 1 and 10")
	}
	if r.Sequence < 1 {
		return fmt.Errorf("sequence must be positive")
	}
	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	return nil
}

func (r Record) IsMutable() bool {
	return r.Status == StatusDraft || r.Status == StatusReviewed
}

func (r Record) CanArchive() bool {
	return r.Status == StatusApproved || r.Status == StatusReviewed
}

func (r Record) Summary() string {
	return fmt.Sprintf("%s:%s/%s intensity=%d", r.Student, r.Gesture, r.Particle, r.Intensity)
}
