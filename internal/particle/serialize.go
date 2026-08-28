package particle

import (
	"encoding/json"
	"gestureparticles/internal/flow021"
)

func EncodeFrame(frame Frame) ([]byte, error) { return json.Marshal(frame) }

func DecodeFrame(data []byte) (Frame, error) {
	var frame Frame
	err := json.Unmarshal(data, &frame)
	return frame, err
}

func FramesFor(observations []flow021.Observation) []Frame {
	result := make([]Frame, 0, len(observations))
	for _, observation := range observations {
		result = append(result, BuildFrame(observation))
	}
	return result
}

func TotalPointCount(frames []Frame) int {
	total := 0
	for _, frame := range frames {
		total += len(frame.Points)
	}
	return total
}
