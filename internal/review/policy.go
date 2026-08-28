package review

import (
	"fmt"
	"gestureparticles/internal/model"
)

type Policy struct {
	RequiredApprovals  int
	AllowedIntensities []int
}

func DefaultPolicy() Policy {
	return Policy{RequiredApprovals: 1, AllowedIntensities: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
}

func (p Policy) Allows(record model.Record) bool {
	for _, intensity := range p.AllowedIntensities {
		if intensity == record.Intensity {
			return true
		}
	}
	return false
}

func (p Policy) Validate(record model.Record) error {
	if !p.Allows(record) {
		return fmt.Errorf("intensity %d is not allowed", record.Intensity)
	}
	if p.RequiredApprovals < 1 {
		return fmt.Errorf("policy requires an approval")
	}
	return nil
}
