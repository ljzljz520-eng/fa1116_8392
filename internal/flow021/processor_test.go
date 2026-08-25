package flow021

import "testing"

func TestProcessorRejectsOutOfOrder(t *testing.T) {
	p := New()
	_, err := p.Process(Observation{ID: "a", Classroom: "room-a", Sequence: 1, Gesture: "wave", Particle: "blue", Content: "one"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Process(Observation{ID: "b", Classroom: "room-a", Sequence: 1, Gesture: "wave", Particle: "blue", Content: "duplicate"}); err == nil {
		t.Fatal("expected sequence error")
	}
}

func TestProcessorHistoryAndSummary(t *testing.T) {
	p := New()
	first, err := p.Process(Observation{ID: "a", Classroom: "room-a", Sequence: 1, Gesture: "wave", Particle: "blue", Content: "one"})
	if err != nil {
		t.Fatal(err)
	}
	if FormatObservation(first) == "" {
		t.Fatal("format empty")
	}
	if p.Count("room-a") != 1 {
		t.Fatal("count mismatch")
	}
	if err := p.ValidateMigration("room-a"); err != nil {
		t.Fatal(err)
	}
}

func Test1116BusinessRegression(t *testing.T) {
	p := New()
	first, err := p.Process(Observation{ID: "gesture-1", Classroom: "room-a", Sequence: 1, Gesture: "wave", Particle: "blue", Content: "raise left hand"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := p.Process(Observation{ID: "gesture-2", Classroom: "room-a", Sequence: 2, Gesture: "pinch", Particle: "gold", Content: "pinch thumb and finger"})
	if err != nil {
		t.Fatal(err)
	}
	if first.Content == second.Content {
		t.Fatalf("second observation reused previous content: %q", second.Content)
	}
}
