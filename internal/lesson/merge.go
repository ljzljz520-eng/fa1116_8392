package lesson

import (
	"fmt"
	"sort"
)

func MergePlans(primary, secondary Plan) (Plan, error) {
	if primary.ID == "" || secondary.ID == "" {
		return Plan{}, fmt.Errorf("both plans are required")
	}
	if primary.Classroom != secondary.Classroom {
		return Plan{}, fmt.Errorf("plans belong to different classrooms")
	}
	merged := primary
	merged.ID = primary.ID + "+" + secondary.ID
	merged.Title = primary.Title + " / " + secondary.Title
	merged.Objective = primary.Objective + "; " + secondary.Objective
	merged.Status = "draft"
	merged.Cues = append([]Cue(nil), primary.Cues...)
	for _, cue := range secondary.Cues {
		cue.ID = secondary.ID + ":" + cue.ID
		if err := merged.AddCue(cue); err != nil {
			return Plan{}, err
		}
	}
	return merged, nil
}

func SortPlans(plans []Plan) []Plan {
	result := append([]Plan(nil), plans...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Classroom == result[j].Classroom {
			return result[i].Title < result[j].Title
		}
		return result[i].Classroom < result[j].Classroom
	})
	return result
}

func PlanIDs(plans []Plan) []string {
	result := make([]string, 0, len(plans))
	for _, plan := range plans {
		result = append(result, plan.ID)
	}
	return result
}

func CountReady(plans []Plan) int {
	count := 0
	for _, plan := range plans {
		if plan.Status == "ready" {
			count++
		}
	}
	return count
}

func TotalCueCount(plans []Plan) int {
	total := 0
	for _, plan := range plans {
		total += len(plan.Cues)
	}
	return total
}

func HasObjective(plans []Plan, objective string) bool {
	for _, plan := range plans {
		if plan.Objective == objective {
			return true
		}
	}
	return false
}

func ClassroomSet(plans []Plan) map[string]bool {
	result := make(map[string]bool)
	for _, plan := range plans {
		result[plan.Classroom] = true
	}
	return result
}

func DraftPlans(plans []Plan) []Plan {
	result := make([]Plan, 0)
	for _, plan := range plans {
		if plan.Status == "draft" {
			result = append(result, plan)
		}
	}
	return result
}

func DeliveredPlans(plans []Plan) []Plan {
	result := make([]Plan, 0)
	for _, plan := range plans {
		if plan.Status == "delivered" {
			result = append(result, plan)
		}
	}
	return result
}

func FilterByTeacher(plans []Plan, teacher string) []Plan {
	result := make([]Plan, 0)
	for _, plan := range plans {
		if plan.Teacher == teacher {
			result = append(result, plan)
		}
	}
	return result
}

func FilterWithLongCue(plans []Plan) []Plan {
	result := make([]Plan, 0)
	for _, plan := range plans {
		if HasLongCue(plan) {
			result = append(result, plan)
		}
	}
	return result
}

func PlanCount(plans []Plan) int { return len(plans) }
