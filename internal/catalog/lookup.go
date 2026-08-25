package catalog

import "strings"

func GestureNames() []string {
	result := make([]string, 0, len(definitions))
	for _, definition := range definitions {
		if definition.Active {
			result = append(result, definition.Name)
		}
	}
	return result
}

func ParticleNames() []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, definition := range definitions {
		for _, particle := range definition.Particles {
			if !seen[particle] {
				seen[particle] = true
				result = append(result, particle)
			}
		}
	}
	return result
}

func Match(query string) []GestureDefinition {
	query = strings.ToLower(strings.TrimSpace(query))
	result := make([]GestureDefinition, 0)
	for _, definition := range definitions {
		if strings.Contains(definition.Name, query) || strings.Contains(definition.Description, query) {
			result = append(result, definition)
		}
	}
	return result
}
