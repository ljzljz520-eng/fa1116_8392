package catalog

import "strings"

type GestureDefinition struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Particles   []string `json:"particles"`
	Active      bool     `json:"active"`
}

var definitions = []GestureDefinition{
	{Name: "wave", Description: "open hand wave", Particles: []string{"blue", "white"}, Active: true},
	{Name: "pinch", Description: "finger pinch", Particles: []string{"gold", "violet"}, Active: true},
	{Name: "point", Description: "directed point", Particles: []string{"green", "amber"}, Active: true},
}

func List() []GestureDefinition { return append([]GestureDefinition(nil), definitions...) }

func Find(name string) (GestureDefinition, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, definition := range definitions {
		if definition.Name == name && definition.Active {
			return definition, true
		}
	}
	return GestureDefinition{}, false
}

func SupportsParticle(gesture, particle string) bool {
	definition, ok := Find(gesture)
	if !ok {
		return false
	}
	for _, allowed := range definition.Particles {
		if strings.EqualFold(allowed, particle) {
			return true
		}
	}
	return false
}
