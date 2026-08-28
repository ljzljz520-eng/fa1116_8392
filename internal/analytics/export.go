package analytics

import (
	"encoding/json"
	"gestureparticles/internal/model"
)

type Export struct {
	Dashboard Dashboard       `json:"dashboard"`
	Timeline  []TimelineEntry `json:"timeline"`
	Insights  []Insight       `json:"insights"`
}

func BuildExport(classroom string, records []model.Record) Export {
	return Export{Dashboard: BuildDashboard(classroom, records), Timeline: BuildTimeline(records), Insights: FindInsights(records)}
}

func MarshalExport(export Export) ([]byte, error) { return json.Marshal(export) }

func FilterApproved(records []model.Record) []model.Record {
	result := make([]model.Record, 0)
	for _, record := range records {
		if record.Status == model.StatusApproved {
			result = append(result, record)
		}
	}
	return result
}
