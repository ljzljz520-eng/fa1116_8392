package store

import "gestureparticles/internal/model"

const workflowBucket = "workflows"

func (d *Database) SaveWorkflow(workflow model.Workflow) error {
	if err := workflow.Validate(); err != nil {
		return err
	}
	return d.Put(workflowBucket, workflow.ID, workflow)
}

func (d *Database) LoadWorkflow(id string) (model.Workflow, error) {
	var workflow model.Workflow
	found, err := d.Get(workflowBucket, id, &workflow)
	if err != nil {
		return workflow, err
	}
	if !found {
		return workflow, &NotFound{Kind: "workflow", ID: id}
	}
	return workflow, nil
}

func (d *Database) ListWorkflows(classroom string) ([]model.Workflow, error) {
	values, err := d.List(workflowBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.Workflow, 0, len(values))
	for _, data := range values {
		var workflow model.Workflow
		if err := unmarshal(data, &workflow); err != nil {
			return nil, err
		}
		if classroom == "" || workflow.Classroom == classroom {
			result = append(result, workflow)
		}
	}
	return result, nil
}

type NotFound struct {
	Kind string
	ID   string
}

func (e *NotFound) Error() string { return e.Kind + " " + e.ID + " not found" }
