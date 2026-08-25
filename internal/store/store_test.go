package store

import (
	"gestureparticles/internal/model"
	"testing"
)

func TestRecordStoreListAndSequence(t *testing.T) {
	db, err := Open(t.TempDir() + "/data.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := model.Record{ID: "r1", Classroom: "room-a", Student: "A", Gesture: "wave", Particle: "blue", Intensity: 3, Sequence: 1, Status: model.StatusDraft, CreatedBy: "teacher"}
	if err := db.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	values, err := db.ListRecords("room-a")
	if err != nil || len(values) != 1 {
		t.Fatalf("values=%v err=%v", values, err)
	}
	next, err := db.NextSequence("room-a")
	if err != nil || next != 2 {
		t.Fatalf("next=%d err=%v", next, err)
	}
}
