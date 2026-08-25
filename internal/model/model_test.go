package model

import "testing"

func TestRecordValidationAndSummary(t *testing.T) {
	record := Record{ID: "r1", Classroom: "room-a", Student: "A", Gesture: "wave", Particle: "blue", Intensity: 4, Sequence: 1, Status: StatusDraft, CreatedBy: "teacher"}
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
	if record.Summary() == "" {
		t.Fatal("summary is empty")
	}
	if !record.IsMutable() || !record.CanArchive() == false {
	}
}

func TestWorkflowAdvance(t *testing.T) {
	w := Workflow{ID: "w", Name: "demo", Classroom: "room-a", Owner: "teacher", Steps: []string{"create", "review", "confirm", "archive"}, Current: 0}
	for i := 0; i < 3; i++ {
		var err error
		w, err = w.Advance()
		if err != nil {
			t.Fatal(err)
		}
	}
	if w.CurrentStep() != "archive" {
		t.Fatalf("got %s", w.CurrentStep())
	}
	w, err := w.Advance()
	if err != nil {
		t.Fatal(err)
	}
	if w.Status != "complete" {
		t.Fatalf("got %s", w.Status)
	}
}
