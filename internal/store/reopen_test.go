package store

import (
	"gestureparticles/internal/model"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persistent.db"
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := model.Record{ID: "persist-1", Classroom: "room-a", Student: "Ada", Gesture: "wave", Particle: "blue", Intensity: 5, Sequence: 1, Status: model.StatusApproved, CreatedBy: "teacher"}
	if err := db.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	loaded, err := db.LoadRecord("persist-1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Student != "Ada" || loaded.Status != model.StatusApproved {
		t.Fatalf("loaded=%+v", loaded)
	}
}
