package analytics

import (
	"gestureparticles/internal/model"
	"testing"
)

func TestDashboardAndTimeline(t *testing.T) {
	records := []model.Record{{ID: "a", Classroom: "room-a", Student: "A", Gesture: "wave", Particle: "blue", Intensity: 2, Sequence: 1, Status: model.StatusApproved}, {ID: "b", Classroom: "room-a", Student: "B", Gesture: "pinch", Particle: "gold", Intensity: 8, Sequence: 3, Status: model.StatusDraft}}
	dashboard := BuildDashboard("room-a", records)
	if dashboard.Total != 2 || dashboard.Approved != 1 {
		t.Fatalf("dashboard=%+v", dashboard)
	}
	if len(MissingSequences(records)) != 1 {
		t.Fatal("missing sequence not detected")
	}
	if TopGesture(records) == "" || len(BuildExport("room-a", records).Insights) != 1 {
		t.Fatal("insight mismatch")
	}
}
