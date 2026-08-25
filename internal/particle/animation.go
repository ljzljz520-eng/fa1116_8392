package particle

import "gestureparticles/internal/flow021"

type Keyframe struct {
	Index      int    `json:"index"`
	Progress   int    `json:"progress"`
	Caption    string `json:"caption"`
	PointCount int    `json:"point_count"`
}

func Keyframes(observation flow021.Observation, count int) []Keyframe {
	if count < 1 {
		return nil
	}
	frame := BuildFrame(observation)
	result := make([]Keyframe, 0, count)
	for index := 0; index < count; index++ {
		result = append(result, Keyframe{Index: index, Progress: index * 100 / (count - 1 + boolInt(count == 1)), Caption: frame.Caption, PointCount: len(frame.Points)})
	}
	return result
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func MergeFrames(frames []Frame) Frame {
	if len(frames) == 0 {
		return Frame{}
	}
	merged := frames[0]
	merged.Points = append([]Point(nil), frames[0].Points...)
	for _, frame := range frames[1:] {
		merged.Points = append(merged.Points, frame.Points...)
	}
	return merged
}
