package catalog

import "fmt"

type RuleSet struct {
	MaxIntensity          int
	RequireParticleMatch  bool
	RequireKnownClassroom bool
}

func DefaultRules() RuleSet {
	return RuleSet{MaxIntensity: 10, RequireParticleMatch: true, RequireKnownClassroom: false}
}

func (r RuleSet) Validate(gesture, particle, classroom string, intensity int) error {
	if intensity < 1 || intensity > r.MaxIntensity {
		return fmt.Errorf("intensity outside allowed range")
	}
	if r.RequireParticleMatch && !SupportsParticle(gesture, particle) {
		return fmt.Errorf("particle is not supported by gesture")
	}
	if r.RequireKnownClassroom {
		if _, ok := FindClassroom(classroom); !ok {
			return fmt.Errorf("unknown classroom")
		}
	}
	return nil
}
