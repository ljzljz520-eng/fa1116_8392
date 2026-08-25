package review

import (
	"gestureparticles/internal/model"
	"sort"
)

func SortByStatus(records []model.Record) []model.Record {
	result := append([]model.Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Status < result[j].Status })
	return result
}

func StatusCounts(records []model.Record) map[model.RecordStatus]int {
	result := make(map[model.RecordStatus]int)
	for _, record := range records {
		result[record.Status]++
	}
	return result
}

func NeedsAttention(record model.Record) bool {
	return record.Status == model.StatusDraft && record.Intensity >= 8
}
