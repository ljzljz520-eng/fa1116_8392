package lesson

import (
	"fmt"
	"sort"
)

type Evaluation struct {
	PlanID   string         `json:"plan_id"`
	Student  string         `json:"student"`
	Scores   map[string]int `json:"scores"`
	Feedback string         `json:"feedback"`
	Complete bool           `json:"complete"`
}

func NewEvaluation(planID, student string) Evaluation {
	return Evaluation{PlanID: planID, Student: student, Scores: make(map[string]int)}
}

func (e *Evaluation) SetScore(cueID string, score int) error {
	if cueID == "" {
		return fmt.Errorf("cue id is required")
	}
	if score < 0 || score > 5 {
		return fmt.Errorf("score must be between zero and five")
	}
	if e.Scores == nil {
		e.Scores = make(map[string]int)
	}
	e.Scores[cueID] = score
	return nil
}

func (e Evaluation) Total() int {
	total := 0
	for _, score := range e.Scores {
		total += score
	}
	return total
}

func (e Evaluation) Average() float64 {
	if len(e.Scores) == 0 {
		return 0
	}
	return float64(e.Total()) / float64(len(e.Scores))
}

func (e *Evaluation) CompleteFor(plan Plan) bool {
	if len(plan.Cues) == 0 {
		return false
	}
	for _, cue := range plan.Cues {
		if _, ok := e.Scores[cue.ID]; !ok {
			return false
		}
	}
	e.Complete = true
	return true
}

func (e Evaluation) RankedCues() []string {
	keys := make([]string, 0, len(e.Scores))
	for key := range e.Scores {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return e.Scores[keys[i]] > e.Scores[keys[j]] })
	return keys
}
