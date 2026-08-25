package lesson

import (
	"fmt"
	"strings"
)

type Cue struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	Gesture         string `json:"gesture"`
	Particle        string `json:"particle"`
	Prompt          string `json:"prompt"`
	ExpectedSeconds int    `json:"expected_seconds"`
}

func (c Cue) Validate() error {
	if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Label) == "" {
		return fmt.Errorf("cue identity is required")
	}
	if strings.TrimSpace(c.Gesture) == "" || strings.TrimSpace(c.Particle) == "" {
		return fmt.Errorf("cue gesture is required")
	}
	if strings.TrimSpace(c.Prompt) == "" {
		return fmt.Errorf("cue prompt is required")
	}
	if c.ExpectedSeconds < 1 || c.ExpectedSeconds > 600 {
		return fmt.Errorf("cue duration is invalid")
	}
	return nil
}

func (c Cue) IsLong() bool { return c.ExpectedSeconds >= 120 }

func (c Cue) Display() string { return fmt.Sprintf("%s: %s (%s)", c.Label, c.Gesture, c.Particle) }

func NormalizeCue(c Cue) Cue {
	c.ID = strings.TrimSpace(c.ID)
	c.Label = strings.TrimSpace(c.Label)
	c.Gesture = strings.ToLower(strings.TrimSpace(c.Gesture))
	c.Particle = strings.ToLower(strings.TrimSpace(c.Particle))
	c.Prompt = strings.TrimSpace(c.Prompt)
	return c
}

func DefaultCues() []Cue {
	return []Cue{{ID: "observe", Label: "Observe", Gesture: "wave", Particle: "blue", Prompt: "Show the wave", ExpectedSeconds: 30}, {ID: "compare", Label: "Compare", Gesture: "pinch", Particle: "gold", Prompt: "Compare particle density", ExpectedSeconds: 60}}
}
