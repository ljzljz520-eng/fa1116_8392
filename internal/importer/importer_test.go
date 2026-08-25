package importer

import (
	"gestureparticles/internal/store"
	"testing"
)

func TestWorkflowImportReport(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/import.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := New(db)
	rows := []Row{{ID: "i1", Classroom: "room-a", Student: "Ada", Gesture: "wave", Particle: "blue", Intensity: "4"}, {ID: "i2", Classroom: "room-a", Student: "Ben", Gesture: "pinch", Particle: "gold", Intensity: "bad"}}
	report, err := s.Run(rows, "teacher")
	if err != nil {
		t.Fatal(err)
	}
	if report.Imported != 1 || report.Rejected != 1 {
		t.Fatalf("report=%+v", report)
	}
	if Render(report) == "" {
		t.Fatal("render empty")
	}
}
