package analytics

import (
	"gestureparticles/internal/model"
	"sort"
)

type TimelineEntry struct {
	Sequence int                `json:"sequence"`
	RecordID string             `json:"record_id"`
	Status   model.RecordStatus `json:"status"`
	Summary  string             `json:"summary"`
}

func BuildTimeline(records []model.Record) []TimelineEntry {
	result := make([]TimelineEntry, 0, len(records))
	for _, record := range records {
		result = append(result, TimelineEntry{Sequence: record.Sequence, RecordID: record.ID, Status: record.Status, Summary: record.Summary()})
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result
}

func MissingSequences(records []model.Record) []int {
	if len(records) == 0 {
		return nil
	}
	timeline := BuildTimeline(records)
	missing := make([]int, 0)
	for expected := 1; expected < timeline[len(timeline)-1].Sequence; expected++ {
		found := false
		for _, item := range timeline {
			if item.Sequence == expected {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, expected)
		}
	}
	return missing
}
