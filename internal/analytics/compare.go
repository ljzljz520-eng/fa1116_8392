package analytics

import "gestureparticles/internal/model"

type Comparison struct {
	Left           string  `json:"left"`
	Right          string  `json:"right"`
	LeftCount      int     `json:"left_count"`
	RightCount     int     `json:"right_count"`
	IntensityDelta float64 `json:"intensity_delta"`
}

func CompareClassrooms(left, right string, records []model.Record) Comparison {
	comparison := Comparison{Left: left, Right: right}
	leftTotal, rightTotal := 0, 0
	for _, record := range records {
		if record.Classroom == left {
			comparison.LeftCount++
			leftTotal += record.Intensity
		}
		if record.Classroom == right {
			comparison.RightCount++
			rightTotal += record.Intensity
		}
	}
	leftAverage, rightAverage := 0.0, 0.0
	if comparison.LeftCount > 0 {
		leftAverage = float64(leftTotal) / float64(comparison.LeftCount)
	}
	if comparison.RightCount > 0 {
		rightAverage = float64(rightTotal) / float64(comparison.RightCount)
	}
	comparison.IntensityDelta = leftAverage - rightAverage
	return comparison
}

func StatusLabel(status model.RecordStatus) string {
	switch status {
	case model.StatusDraft:
		return "Needs review"
	case model.StatusReviewed:
		return "Reviewed"
	case model.StatusApproved:
		return "Approved"
	case model.StatusArchived:
		return "Archived"
	default:
		return "Unknown"
	}
}

func CountDistinctStudents(records []model.Record) int {
	students := make(map[string]bool)
	for _, record := range records {
		students[record.Student] = true
	}
	return len(students)
}

func CountDistinctParticles(records []model.Record) int {
	particles := make(map[string]bool)
	for _, record := range records {
		particles[record.Particle] = true
	}
	return len(particles)
}

func HighIntensity(records []model.Record, threshold int) []model.Record {
	result := make([]model.Record, 0)
	for _, record := range records {
		if record.Intensity >= threshold {
			result = append(result, record)
		}
	}
	return result
}
