package lesson

import (
	"fmt"
	"sort"
)

type Slot struct {
	PlanID    string `json:"plan_id"`
	Classroom string `json:"classroom"`
	Order     int    `json:"order"`
	Label     string `json:"label"`
}

func BuildSchedule(plan Plan) ([]Slot, error) {
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	slots := make([]Slot, 0, len(plan.Cues))
	for index, cue := range plan.Cues {
		slots = append(slots, Slot{PlanID: plan.ID, Classroom: plan.Classroom, Order: index + 1, Label: cue.Display()})
	}
	sort.SliceStable(slots, func(i, j int) bool { return slots[i].Order < slots[j].Order })
	return slots, nil
}

func FindSlot(slots []Slot, order int) (Slot, error) {
	for _, slot := range slots {
		if slot.Order == order {
			return slot, nil
		}
	}
	return Slot{}, fmt.Errorf("slot %d not found", order)
}

func Duration(plan Plan) int {
	total := 0
	for _, cue := range plan.Cues {
		total += cue.ExpectedSeconds
	}
	return total
}

func HasLongCue(plan Plan) bool {
	for _, cue := range plan.Cues {
		if cue.IsLong() {
			return true
		}
	}
	return false
}
