package analytics

import (
	"fmt"
	"gestureparticles/internal/model"
)

type Insight struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

func FindInsights(records []model.Record) []Insight {
	insights := make([]Insight, 0)
	for _, record := range records {
		if record.Status == model.StatusDraft && record.Intensity >= 8 {
			insights = append(insights, Insight{Code: "HIGH_DRAFT_INTENSITY", Message: fmt.Sprintf("%s needs review", record.Student), Severity: "attention"})
		}
		if record.Status == model.StatusArchived && record.Version < 2 {
			insights = append(insights, Insight{Code: "UNREVIEWED_ARCHIVE", Message: record.ID, Severity: "warning"})
		}
	}
	return insights
}

func TopGesture(records []model.Record) string {
	counts := make(map[string]int)
	best := ""
	bestCount := 0
	for _, record := range records {
		counts[record.Gesture]++
		if counts[record.Gesture] > bestCount {
			best = record.Gesture
			bestCount = counts[record.Gesture]
		}
	}
	return best
}
