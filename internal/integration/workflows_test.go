package integration

import (
	"gestureparticles/internal/flow021"
	"gestureparticles/internal/importer"
	"gestureparticles/internal/model"
	"gestureparticles/internal/registry"
	"gestureparticles/internal/review"
	"gestureparticles/internal/store"
	"testing"
)

func newService(t *testing.T) (*registry.Service, *store.Database) {
	db, err := store.Open(t.TempDir() + "/workflow.db")
	if err != nil {
		t.Fatal(err)
	}
	return registry.New(db, review.New(db), importer.New(db), flow021.New()), db
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, db := newService(t)
	defer db.Close()
	record, err := s.Register(model.Record{ID: "wf-1", Classroom: "room-a", Student: "Ada", Gesture: "wave", Particle: "blue", Intensity: 5, CreatedBy: "teacher"})
	if err != nil {
		t.Fatal(err)
	}
	record, err = s.Review(record.ID, "reviewer", "approve", "ready")
	if err != nil || record.Status != model.StatusApproved {
		t.Fatalf("review=%+v err=%v", record, err)
	}
	record, err = s.Archive(record.ID, "archiver")
	if err != nil || record.Status != model.StatusArchived {
		t.Fatalf("archive=%+v err=%v", record, err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, db := newService(t)
	defer db.Close()
	_, err := s.Register(model.Record{ID: "wf-2", Classroom: "room-a", Student: "Ben", Gesture: "pinch", Particle: "gold", Intensity: 4, CreatedBy: "teacher"})
	if err != nil {
		t.Fatal(err)
	}
	records, err := s.Search("room-a", "ben", model.StatusDraft)
	if err != nil || len(records) != 1 {
		t.Fatalf("records=%v err=%v", records, err)
	}
	updated, err := s.Update("wf-2", model.Record{Notes: "published candidate"}, "teacher")
	if err != nil || updated.Notes != "published candidate" {
		t.Fatalf("updated=%+v err=%v", updated, err)
	}
	if _, err := s.Review("wf-2", "reviewer", "approve", "publish"); err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	s, db := newService(t)
	defer db.Close()
	report, err := s.Import([]importer.Row{{ID: "wf-3", Classroom: "room-b", Student: "Kai", Gesture: "point", Particle: "green", Intensity: "6"}}, "teacher")
	if err != nil || report.Imported != 1 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	snapshot, err := s.Snapshot("room-b")
	if err != nil || snapshot.Dashboard.Total != 1 {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
}
