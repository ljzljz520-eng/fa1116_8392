package lesson

import "strings"

type Summary struct {
	PlanID          string   `json:"plan_id"`
	Title           string   `json:"title"`
	CueLabels       []string `json:"cue_labels"`
	DurationSeconds int      `json:"duration_seconds"`
	Ready           bool     `json:"ready"`
}

func BuildSummary(plan Plan) Summary {
	labels := make([]string, 0, len(plan.Cues))
	for _, cue := range plan.Cues {
		labels = append(labels, cue.Label)
	}
	return Summary{PlanID: plan.ID, Title: plan.Title, CueLabels: labels, DurationSeconds: Duration(plan), Ready: plan.Status == "ready" || plan.Status == "delivered"}
}

func SearchPlans(plans []Plan, query string) []Plan {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return append([]Plan(nil), plans...)
	}
	result := make([]Plan, 0)
	for _, plan := range plans {
		if strings.Contains(strings.ToLower(plan.Title+" "+plan.Objective+" "+plan.Classroom), query) {
			result = append(result, plan)
		}
	}
	return result
}

func ReadyPlans(plans []Plan) []Plan {
	result := make([]Plan, 0)
	for _, plan := range plans {
		if plan.Status == "ready" {
			result = append(result, plan)
		}
	}
	return result
}
