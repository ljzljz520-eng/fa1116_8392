package importer

import (
	"gestureparticles/internal/model"
	"strings"
)

func Render(report Report) string {
	lines := []string{"imported=" + itoa(report.Imported), "rejected=" + itoa(report.Rejected)}
	for _, err := range report.Errors {
		lines = append(lines, "error="+err)
	}
	for _, record := range report.Records {
		lines = append(lines, record.Summary())
	}
	return strings.Join(lines, "\n")
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	digits := make([]byte, 0, 12)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}

func GroupByClassroom(records []model.Record) map[string][]model.Record {
	result := make(map[string][]model.Record)
	for _, record := range records {
		result[record.Classroom] = append(result[record.Classroom], record)
	}
	return result
}
