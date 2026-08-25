package review

import (
	"fmt"
	"gestureparticles/internal/model"
	"gestureparticles/internal/store"
	"strings"
)

type Service struct{ db *store.Database }

func New(db *store.Database) *Service { return &Service{db: db} }

func (s *Service) Decide(record model.Record, actor, decision, reason string) (model.Record, error) {
	actor = strings.TrimSpace(actor)
	decision = strings.ToLower(strings.TrimSpace(decision))
	if actor == "" {
		return record, fmt.Errorf("reviewer is required")
	}
	if decision != "approve" && decision != "reject" {
		return record, fmt.Errorf("decision must be approve or reject")
	}
	if record.Status != model.StatusDraft && record.Status != model.StatusReviewed {
		return record, fmt.Errorf("record is not reviewable")
	}
	old := record.Status
	record.Reviewer = actor
	record.ReviewedAt = "deterministic-review"
	if decision == "approve" {
		record.Status = model.StatusApproved
	} else {
		record.Status = model.StatusDraft
		record.Notes = strings.TrimSpace(record.Notes + " " + reason)
	}
	record.Version++
	if err := s.RecordTransition(record, "review", actor, string(old), string(record.Status)); err != nil {
		return record, err
	}
	return record, nil
}

func (s *Service) RecordTransition(record model.Record, action, actor, from, to string) error {
	event := model.AuditEvent{ID: record.ID + ":" + action + ":" + fmt.Sprint(record.Version), RecordID: record.ID, Action: action, Actor: actor, Reason: record.Notes, Timestamp: "deterministic", From: from, To: to}
	return s.db.SaveAudit(event)
}

func (s *Service) CanReview(record model.Record) bool {
	return record.Status == model.StatusDraft || record.Status == model.StatusReviewed
}

func (s *Service) Pending(classroom string) ([]model.Record, error) {
	records, err := s.db.ListRecords(classroom)
	if err != nil {
		return nil, err
	}
	result := make([]model.Record, 0)
	for _, record := range records {
		if s.CanReview(record) {
			result = append(result, record)
		}
	}
	return result, nil
}

func (s *Service) AuditSummary(recordID string) (map[string]int, error) {
	events, err := s.db.ListAudits(recordID)
	if err != nil {
		return nil, err
	}
	result := make(map[string]int)
	for _, event := range events {
		result[event.Action]++
	}
	return result, nil
}
