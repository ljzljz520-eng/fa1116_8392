package analytics

import (
	"gestureparticles/internal/model"
	"sort"
)

type Dashboard struct {
	Classroom        string   `json:"classroom"`
	Total            int      `json:"total"`
	Draft            int      `json:"draft"`
	Reviewed         int      `json:"reviewed"`
	Approved         int      `json:"approved"`
	Archived         int      `json:"archived"`
	AverageIntensity float64  `json:"average_intensity"`
	Gestures         []string `json:"gestures"`
}

func BuildDashboard(classroom string, records []model.Record) Dashboard {
	dashboard := Dashboard{Classroom: classroom}
	gestures := make(map[string]bool)
	totalIntensity := 0
	for _, record := range records {
		dashboard.Total++
		totalIntensity += record.Intensity
		gestures[record.Gesture] = true
		switch record.Status {
		case model.StatusDraft:
			dashboard.Draft++
		case model.StatusReviewed:
			dashboard.Reviewed++
		case model.StatusApproved:
			dashboard.Approved++
		case model.StatusArchived:
			dashboard.Archived++
		}
	}
	if dashboard.Total > 0 {
		dashboard.AverageIntensity = float64(totalIntensity) / float64(dashboard.Total)
	}
	for gesture := range gestures {
		dashboard.Gestures = append(dashboard.Gestures, gesture)
	}
	sort.Strings(dashboard.Gestures)
	return dashboard
}

func CompletionRatio(dashboard Dashboard) float64 {
	if dashboard.Total == 0 {
		return 0
	}
	return float64(dashboard.Approved+dashboard.Archived) / float64(dashboard.Total)
}
