package review

import (
	"gestureparticles/internal/model"
	"gestureparticles/internal/store"
	"testing"
)

func TestReviewApproveAndPending(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/review.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := New(db)
	r := model.Record{ID: "r1", Classroom: "room-a", Student: "A", Gesture: "wave", Particle: "blue", Intensity: 3, Sequence: 1, Status: model.StatusDraft, CreatedBy: "teacher"}
	updated, err := s.Decide(r, "reviewer", "approve", "ok")
	if err != nil || updated.Status != model.StatusApproved {
		t.Fatalf("updated=%+v err=%v", updated, err)
	}
	events, err := db.ListAudits("r1")
	if err != nil || len(events) != 1 {
		t.Fatalf("events=%v err=%v", events, err)
	}
}
