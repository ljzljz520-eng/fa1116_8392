package registry

import (
	"fmt"
	"gestureparticles/internal/model"
)

func (s *Service) StartWorkflow(workflow model.Workflow) (model.Workflow, error) {
	if workflow.Status == "" {
		workflow.Status = "active"
	}
	if err := workflow.Validate(); err != nil {
		return workflow, err
	}
	if err := s.db.SaveWorkflow(workflow); err != nil {
		return workflow, err
	}
	return workflow, nil
}

func (s *Service) AdvanceWorkflow(id string) (model.Workflow, error) {
	workflow, err := s.db.LoadWorkflow(id)
	if err != nil {
		return workflow, err
	}
	if workflow.Status == "complete" {
		return workflow, fmt.Errorf("workflow %s is complete", id)
	}
	updated, err := workflow.Advance()
	if err != nil {
		return workflow, err
	}
	if err := s.db.SaveWorkflow(updated); err != nil {
		return workflow, err
	}
	return updated, nil
}

func (s *Service) GetWorkflow(id string) (model.Workflow, error) { return s.db.LoadWorkflow(id) }
