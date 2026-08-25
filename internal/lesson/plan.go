package lesson

import (
	"fmt"
	"strings"
)

type Plan struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Classroom string `json:"classroom"`
	Teacher   string `json:"teacher"`
	Objective string `json:"objective"`
	Cues      []Cue  `json:"cues"`
	Status    string `json:"status"`
}

func NewPlan(id, title, classroom, teacher, objective string) Plan {
	return Plan{ID: id, Title: title, Classroom: classroom, Teacher: teacher, Objective: objective, Status: "draft", Cues: make([]Cue, 0)}
}

func (p Plan) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Title) == "" {
		return fmt.Errorf("lesson identity is required")
	}
	if strings.TrimSpace(p.Classroom) == "" || strings.TrimSpace(p.Teacher) == "" {
		return fmt.Errorf("lesson ownership is required")
	}
	if strings.TrimSpace(p.Objective) == "" {
		return fmt.Errorf("lesson objective is required")
	}
	if len(p.Cues) < 2 {
		return fmt.Errorf("lesson needs at least two cues")
	}
	if p.Status != "draft" && p.Status != "ready" && p.Status != "delivered" {
		return fmt.Errorf("unknown lesson status")
	}
	return nil
}

func (p *Plan) AddCue(cue Cue) error {
	if err := cue.Validate(); err != nil {
		return err
	}
	for _, existing := range p.Cues {
		if existing.ID == cue.ID {
			return fmt.Errorf("cue %s already exists", cue.ID)
		}
	}
	p.Cues = append(p.Cues, cue)
	return nil
}

func (p *Plan) MarkReady() error {
	if len(p.Cues) < 2 {
		return fmt.Errorf("lesson needs cues before ready")
	}
	if p.Status == "delivered" {
		return fmt.Errorf("delivered lesson cannot return to ready")
	}
	p.Status = "ready"
	return nil
}

func (p *Plan) MarkDelivered() error {
	if p.Status != "ready" {
		return fmt.Errorf("lesson must be ready before delivery")
	}
	p.Status = "delivered"
	return nil
}

func (p Plan) CueCount() int { return len(p.Cues) }

func (p Plan) CueAt(index int) (Cue, bool) {
	if index < 0 || index >= len(p.Cues) {
		return Cue{}, false
	}
	return p.Cues[index], true
}
