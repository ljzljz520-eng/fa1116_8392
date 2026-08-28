package registry

import (
	"errors"
	"fmt"
	"gestureparticles/internal/analytics"
	"gestureparticles/internal/flow021"
	"gestureparticles/internal/importer"
	"gestureparticles/internal/lesson"
	"gestureparticles/internal/model"
	"gestureparticles/internal/review"
	"gestureparticles/internal/store"
	"strings"
	"sync"
)

type Service struct {
	db        *store.Database
	reviews   *review.Service
	imports   *importer.Service
	particles *flow021.Processor
	mu        sync.Mutex
}

func New(db *store.Database, reviews *review.Service, imports *importer.Service, particles *flow021.Processor) *Service {
	return &Service{db: db, reviews: reviews, imports: imports, particles: particles}
}

func (s *Service) Register(input model.Record) (model.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if input.Status == "" {
		input.Status = model.StatusDraft
	}
	if input.Version == 0 {
		input.Version = 1
	}
	if input.Sequence == 0 {
		sequence, err := s.db.NextSequence(input.Classroom)
		if err != nil {
			return input, err
		}
		input.Sequence = sequence
	}
	if err := input.Validate(); err != nil {
		return input, err
	}
	if existing, err := s.db.LoadRecord(input.ID); err == nil && existing.ID != "" {
		return input, fmt.Errorf("record %s already exists", input.ID)
	}
	if err := s.db.SaveRecord(input); err != nil {
		return input, err
	}
	if err := s.reviews.RecordTransition(input, "register", input.CreatedBy, "", string(input.Status)); err != nil {
		return input, err
	}
	return input, nil
}

func (s *Service) Get(id string) (model.Record, error) {
	if strings.TrimSpace(id) == "" {
		return model.Record{}, errors.New("record id is required")
	}
	return s.db.LoadRecord(id)
}

func (s *Service) Search(classroom, query string, status model.RecordStatus) ([]model.Record, error) {
	records, err := s.db.ListRecords(classroom)
	if err != nil {
		return nil, err
	}
	query = strings.ToLower(strings.TrimSpace(query))
	filtered := make([]model.Record, 0, len(records))
	for _, record := range records {
		if status != "" && record.Status != status {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(record.Student+" "+record.Gesture+" "+record.Particle+" "+record.Notes), query) {
			continue
		}
		filtered = append(filtered, record)
	}
	return filtered, nil
}

func (s *Service) Update(id string, patch model.Record, actor string) (model.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, err := s.db.LoadRecord(id)
	if err != nil {
		return current, err
	}
	if !current.IsMutable() {
		return current, fmt.Errorf("record %s cannot be changed while %s", id, current.Status)
	}
	if patch.Student != "" {
		current.Student = patch.Student
	}
	if patch.Gesture != "" {
		current.Gesture = patch.Gesture
	}
	if patch.Particle != "" {
		current.Particle = patch.Particle
	}
	if patch.Intensity != 0 {
		current.Intensity = patch.Intensity
	}
	if patch.Notes != "" {
		current.Notes = patch.Notes
	}
	current.Version++
	if err := current.Validate(); err != nil {
		return current, err
	}
	if err := s.db.SaveRecord(current); err != nil {
		return current, err
	}
	if err := s.reviews.RecordTransition(current, "update", actor, string(current.Status), string(current.Status)); err != nil {
		return current, err
	}
	return current, nil
}

func (s *Service) Review(id, actor, decision, reason string) (model.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, err := s.db.LoadRecord(id)
	if err != nil {
		return record, err
	}
	updated, err := s.reviews.Decide(record, actor, decision, reason)
	if err != nil {
		return record, err
	}
	if err := s.db.SaveRecord(updated); err != nil {
		return record, err
	}
	return updated, nil
}

func (s *Service) Archive(id, actor string) (model.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, err := s.db.LoadRecord(id)
	if err != nil {
		return record, err
	}
	if !record.CanArchive() {
		return record, fmt.Errorf("record %s is not ready to archive", id)
	}
	old := record.Status
	record.Status = model.StatusArchived
	record.Version++
	if err := s.db.SaveRecord(record); err != nil {
		return record, err
	}
	if err := s.reviews.RecordTransition(record, "archive", actor, string(old), string(record.Status)); err != nil {
		return record, err
	}
	return record, nil
}

func (s *Service) AddObservation(input flow021.Observation) (flow021.Observation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.particles.Process(input)
}

func (s *Service) Import(rows []importer.Row, actor string) (importer.Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.imports.Run(rows, actor)
}

func (s *Service) Attach(attachment model.Attachment) error {
	if err := attachment.Validate(); err != nil {
		return err
	}
	if _, err := s.db.LoadRecord(attachment.RecordID); err != nil {
		return err
	}
	return s.db.SaveAttachment(attachment)
}

func (s *Service) History(id string) ([]model.AuditEvent, error) { return s.db.ListAudits(id) }

func (s *Service) PrepareLesson(plan lesson.Plan) ([]lesson.Slot, error) {
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	return lesson.BuildSchedule(plan)
}

func (s *Service) EvaluateLesson(plan lesson.Plan, evaluation lesson.Evaluation) (lesson.Evaluation, error) {
	if err := plan.Validate(); err != nil {
		return evaluation, err
	}
	if evaluation.PlanID != plan.ID {
		return evaluation, fmt.Errorf("evaluation belongs to another plan")
	}
	if !evaluation.CompleteFor(plan) {
		return evaluation, fmt.Errorf("evaluation is missing cue scores")
	}
	return evaluation, nil
}

func (s *Service) Snapshot(classroom string) (Snapshot, error) {
	records, err := s.Search(classroom, "", "")
	if err != nil {
		return Snapshot{}, err
	}
	workflows, err := s.db.ListWorkflows(classroom)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Classroom: classroom, Records: records, Workflows: workflows, Dashboard: analytics.BuildDashboard(classroom, records), Insights: analytics.FindInsights(records)}, nil
}

type Snapshot struct {
	Classroom string              `json:"classroom"`
	Records   []model.Record      `json:"records"`
	Workflows []model.Workflow    `json:"workflows"`
	Dashboard analytics.Dashboard `json:"dashboard"`
	Insights  []analytics.Insight `json:"insights"`
}
