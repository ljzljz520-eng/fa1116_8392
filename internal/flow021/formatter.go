package flow021

import "fmt"

func FormatObservation(value Observation) string {
	return fmt.Sprintf("%s #%d %s %s: %s", value.Classroom, value.Sequence, value.Gesture, value.Particle, value.Content)
}

func FormatTransition(value Transition) string {
	if value.From == "" {
		return "start -> " + value.To
	}
	return value.From + " -> " + value.To
}

func Summarize(values []Observation) map[string]int {
	result := make(map[string]int)
	for _, value := range values {
		result[value.Gesture]++
	}
	return result
}
