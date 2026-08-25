package particle

import "strings"

type Color struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

var palette = map[string]Color{"blue": {Name: "blue", Hex: "#3B82F6"}, "white": {Name: "white", Hex: "#F8FAFC"}, "gold": {Name: "gold", Hex: "#F59E0B"}, "violet": {Name: "violet", Hex: "#8B5CF6"}, "green": {Name: "green", Hex: "#10B981"}, "amber": {Name: "amber", Hex: "#FBBF24"}}

func Resolve(name string) Color {
	if color, ok := palette[strings.ToLower(strings.TrimSpace(name))]; ok {
		return color
	}
	return Color{Name: "neutral", Hex: "#94A3B8"}
}

func Colors() []Color {
	result := make([]Color, 0, len(palette))
	for _, color := range palette {
		result = append(result, color)
	}
	return result
}

func Contrast(color Color) string {
	if color.Name == "white" || color.Name == "gold" || color.Name == "amber" {
		return "dark"
	}
	return "light"
}
