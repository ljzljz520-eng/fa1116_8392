package model

import "fmt"

type Workflow struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Classroom string   `json:"classroom"`
	Owner     string   `json:"owner"`
	Steps     []string `json:"steps"`
	Current   int      `json:"current"`
	Status    string   `json:"status"`
}

func (w Workflow) Validate() error {
	if w.ID == "" || w.Name == "" || w.Classroom == "" || w.Owner == "" {
		return fmt.Errorf("workflow identity is required")
	}
	if len(w.Steps) < 4 {
		return fmt.Errorf("workflow needs at least four steps")
	}
	if w.Current < 0 || w.Current >= len(w.Steps) {
		return fmt.Errorf("workflow step is out of range")
	}
	return nil
}

func (w Workflow) CurrentStep() string {
	if len(w.Steps) == 0 || w.Current >= len(w.Steps) {
		return ""
	}
	return w.Steps[w.Current]
}

func (w Workflow) Advance() (Workflow, error) {
	if err := w.Validate(); err != nil {
		return w, err
	}
	if w.Current+1 >= len(w.Steps) {
		w.Status = "complete"
		return w, nil
	}
	w.Current++
	w.Status = "active"
	return w, nil
}
