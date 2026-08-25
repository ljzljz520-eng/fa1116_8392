package particle

import "gestureparticles/internal/flow021"

type Sequence struct {
	Classroom   string  `json:"classroom"`
	Frames      []Frame `json:"frames"`
	TotalPoints int     `json:"total_points"`
}

func BuildSequence(observations []flow021.Observation) Sequence {
	sequence := Sequence{Frames: make([]Frame, 0, len(observations))}
	for _, observation := range observations {
		if sequence.Classroom == "" {
			sequence.Classroom = observation.Classroom
		}
		frame := BuildFrame(observation)
		sequence.TotalPoints += len(frame.Points)
		sequence.Frames = append(sequence.Frames, frame)
	}
	return sequence
}

func (s Sequence) FrameAt(index int) (Frame, bool) {
	if index < 0 || index >= len(s.Frames) {
		return Frame{}, false
	}
	return s.Frames[index], true
}

func (s Sequence) Empty() bool { return len(s.Frames) == 0 }

func (s Sequence) Captions() []string {
	result := make([]string, 0, len(s.Frames))
	for _, frame := range s.Frames {
		result = append(result, frame.Caption)
	}
	return result
}
