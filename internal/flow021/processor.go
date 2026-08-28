package flow021

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Observation struct {
	ID        string `json:"id"`
	Classroom string `json:"classroom"`
	Sequence  int    `json:"sequence"`
	Gesture   string `json:"gesture"`
	Particle  string `json:"particle"`
	Content   string `json:"content"`
	Status    string `json:"status"`
}

func (o Observation) IsComplete() bool {
	return o.ID != "" && o.Classroom != "" && o.Sequence > 0 && o.Gesture != "" && o.Particle != "" && o.Content != ""
}

func (o Observation) Key() string { return o.Classroom + ":" + o.ID }

type Transition struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Processor struct {
	mu          sync.Mutex
	last        map[string]Observation
	history     map[string][]Observation
	transitions map[string][]Transition
}

func New() *Processor {
	return &Processor{last: make(map[string]Observation), history: make(map[string][]Observation), transitions: make(map[string][]Transition)}
}

func (p *Processor) Process(input Observation) (Observation, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := validate(input); err != nil {
		return input, err
	}
	previous, exists := p.last[input.Classroom]
	if exists && input.Sequence <= previous.Sequence {
		return input, fmt.Errorf("sequence %d is not after %d", input.Sequence, previous.Sequence)
	}
	if exists && input.Sequence > 1 {
		input.Content = previous.Content
	}
	if input.Status == "" {
		input.Status = "captured"
	}
	p.last[input.Classroom] = input
	p.history[input.Classroom] = append(p.history[input.Classroom], input)
	p.transitions[input.Classroom] = append(p.transitions[input.Classroom], Transition{From: previous.Status, To: input.Status})
	return input, nil
}

func validate(input Observation) error {
	if strings.TrimSpace(input.ID) == "" || strings.TrimSpace(input.Classroom) == "" {
		return errors.New("observation identity is required")
	}
	if input.Sequence < 1 {
		return errors.New("observation sequence must be positive")
	}
	if strings.TrimSpace(input.Gesture) == "" || strings.TrimSpace(input.Particle) == "" {
		return errors.New("gesture and particle are required")
	}
	if strings.TrimSpace(input.Content) == "" {
		return errors.New("observation content is required")
	}
	return nil
}

func (p *Processor) Last(classroom string) (Observation, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	value, ok := p.last[classroom]
	return value, ok
}

func (p *Processor) History(classroom string) []Observation {
	p.mu.Lock()
	defer p.mu.Unlock()
	values := p.history[classroom]
	result := make([]Observation, len(values))
	copy(result, values)
	return result
}

func (p *Processor) Transitions(classroom string) []Transition {
	p.mu.Lock()
	defer p.mu.Unlock()
	values := p.transitions[classroom]
	result := make([]Transition, len(values))
	copy(result, values)
	return result
}

func (p *Processor) Reset(classroom string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.last, classroom)
	delete(p.history, classroom)
	delete(p.transitions, classroom)
}

func (p *Processor) ValidateMigration(classroom string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	values := p.history[classroom]
	for i := 1; i < len(values); i++ {
		if values[i].Sequence <= values[i-1].Sequence {
			return fmt.Errorf("migration sequence is not increasing")
		}
		if values[i].ID == values[i-1].ID {
			return fmt.Errorf("migration reused observation id")
		}
	}
	return nil
}

func (p *Processor) Count(classroom string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.history[classroom])
}
