package particle

import (
	"fmt"
	"gestureparticles/internal/flow021"
)

type Point struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Alpha int `json:"alpha"`
}

type Frame struct {
	ObservationID string  `json:"observation_id"`
	Classroom     string  `json:"classroom"`
	Gesture       string  `json:"gesture"`
	Particle      string  `json:"particle"`
	Points        []Point `json:"points"`
	Caption       string  `json:"caption"`
}

func BuildFrame(observation flow021.Observation) Frame {
	points := make([]Point, 0, len(observation.Content)+observation.Sequence)
	seed := deterministicSeed(observation.ID + observation.Content)
	for index := 0; index < len(observation.Content)+observation.Sequence; index++ {
		seed = next(seed)
		points = append(points, Point{X: seed % 640, Y: (seed / 640) % 360, Alpha: 40 + seed%216})
	}
	return Frame{ObservationID: observation.ID, Classroom: observation.Classroom, Gesture: observation.Gesture, Particle: observation.Particle, Points: points, Caption: fmt.Sprintf("%s #%d: %s", observation.Gesture, observation.Sequence, observation.Content)}
}

func deterministicSeed(value string) int {
	seed := 17
	for _, char := range value {
		seed = (seed*31 + int(char)) % 1000003
	}
	return seed
}

func next(seed int) int { return (seed*1103515245 + 12345) & 0x7fffffff }

func Bounds(frame Frame) (int, int) {
	maxX, maxY := 0, 0
	for _, point := range frame.Points {
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}
	return maxX, maxY
}
