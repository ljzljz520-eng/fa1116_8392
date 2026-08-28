package catalog

import "strings"

type Classroom struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Teacher string `json:"teacher"`
	Active  bool   `json:"active"`
}

var classrooms = []Classroom{{ID: "room-a", Name: "Aurora Lab", Teacher: "Ms. Chen", Active: true}, {ID: "room-b", Name: "Orbit Lab", Teacher: "Mr. Lin", Active: true}}

func ListClassrooms() []Classroom { return append([]Classroom(nil), classrooms...) }

func FindClassroom(id string) (Classroom, bool) {
	for _, classroom := range classrooms {
		if classroom.ID == strings.TrimSpace(id) && classroom.Active {
			return classroom, true
		}
	}
	return Classroom{}, false
}

func ClassroomName(id string) string {
	classroom, ok := FindClassroom(id)
	if !ok {
		return ""
	}
	return classroom.Name
}
