package importer

import (
	"fmt"
	"strings"
)

func ValidateRows(rows []Row) []string {
	errors := make([]string, 0)
	seen := make(map[string]bool)
	for i, row := range rows {
		if strings.TrimSpace(row.ID) != "" && seen[row.ID] {
			errors = append(errors, fmt.Sprintf("row %d duplicates %s", i+1, row.ID))
		}
		if row.ID != "" {
			seen[row.ID] = true
		}
		if strings.TrimSpace(row.Classroom) == "" {
			errors = append(errors, fmt.Sprintf("row %d has no classroom", i+1))
		}
	}
	return errors
}

func Normalize(rows []Row) []Row {
	result := make([]Row, len(rows))
	for i, row := range rows {
		row.ID = strings.TrimSpace(row.ID)
		row.Classroom = strings.TrimSpace(row.Classroom)
		row.Student = strings.TrimSpace(row.Student)
		row.Gesture = strings.TrimSpace(row.Gesture)
		row.Particle = strings.TrimSpace(row.Particle)
		row.Intensity = strings.TrimSpace(row.Intensity)
		result[i] = row
	}
	return result
}
