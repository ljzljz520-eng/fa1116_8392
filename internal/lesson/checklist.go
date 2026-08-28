package lesson

import "fmt"

type ChecklistItem struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Done     bool   `json:"done"`
}

func BuildChecklist(plan Plan) []ChecklistItem {
	items := []ChecklistItem{{Name: "objective stated", Required: true, Done: plan.Objective != ""}, {Name: "classroom assigned", Required: true, Done: plan.Classroom != ""}, {Name: "teacher assigned", Required: true, Done: plan.Teacher != ""}, {Name: "cues prepared", Required: true, Done: len(plan.Cues) >= 2}}
	return items
}

func ChecklistComplete(items []ChecklistItem) bool {
	for _, item := range items {
		if item.Required && !item.Done {
			return false
		}
	}
	return true
}

func ChecklistMessage(items []ChecklistItem) string {
	missing := ""
	for _, item := range items {
		if item.Required && !item.Done {
			if missing != "" {
				missing += ", "
			}
			missing += item.Name
		}
	}
	if missing == "" {
		return "ready"
	}
	return fmt.Sprintf("missing: %s", missing)
}
