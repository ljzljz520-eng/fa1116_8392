package importer

import (
	"fmt"
	"gestureparticles/internal/catalog"
	"gestureparticles/internal/model"
	"gestureparticles/internal/store"
	"strconv"
	"strings"
)

type Row struct {
	ID        string
	Classroom string
	Student   string
	Gesture   string
	Particle  string
	Intensity string
	Notes     string
}

type Report struct {
	Imported int            `json:"imported"`
	Rejected int            `json:"rejected"`
	Errors   []string       `json:"errors,omitempty"`
	Records  []model.Record `json:"records"`
}

type Service struct{ db *store.Database }

func New(db *store.Database) *Service { return &Service{db: db} }

func (s *Service) Run(rows []Row, actor string) (Report, error) {
	report := Report{Records: make([]model.Record, 0, len(rows))}
	for index, row := range rows {
		record, err := Parse(row, index, actor)
		if err != nil {
			report.Rejected++
			report.Errors = append(report.Errors, fmt.Sprintf("row %d: %v", index+1, err))
			continue
		}
		if err := s.db.SaveRecord(record); err != nil {
			return report, err
		}
		report.Imported++
		report.Records = append(report.Records, record)
	}
	return report, nil
}

func Parse(row Row, index int, actor string) (model.Record, error) {
	if strings.TrimSpace(row.ID) == "" {
		row.ID = fmt.Sprintf("import-%03d", index+1)
	}
	intensity, err := strconv.Atoi(strings.TrimSpace(row.Intensity))
	if err != nil {
		return model.Record{}, fmt.Errorf("intensity is not numeric")
	}
	record := model.Record{ID: row.ID, Classroom: row.Classroom, Student: row.Student, Gesture: row.Gesture, Particle: row.Particle, Intensity: intensity, Status: model.StatusDraft, CreatedBy: actor, Sequence: index + 1, Version: 1, Notes: row.Notes}
	if err := record.Validate(); err != nil {
		return model.Record{}, err
	}
	if err := catalog.DefaultRules().Validate(record.Gesture, record.Particle, record.Classroom, record.Intensity); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func ParseLines(lines []string, actor string) []Row {
	rows := make([]Row, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) < 6 {
			continue
		}
		rows = append(rows, Row{ID: parts[0], Classroom: parts[1], Student: parts[2], Gesture: parts[3], Particle: parts[4], Intensity: parts[5], Notes: strings.Join(parts[6:], "|")})
	}
	return rows
}
