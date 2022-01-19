package course

import (
	"strings"
)

var Courses map[string] *Course

type Course struct {
	CourseCode string
	CourseColor int
}

func InitializeCourses() {
	Courses = make(map[string] *Course);
}

func CreateCourse(courseCode string, courseColor int) *Course {
	course := Course {
		CourseCode: courseCode,
		CourseColor: courseColor,
	}

	return &course;
}

func ParseCourseCode(summary string) string {
	courseMarkers := []string { "H1", "Y1" };

	// Default to the entire string if no courseMarker found
	index := len(summary) - 1;
	for _, courseMarker := range courseMarkers {
		index = strings.Index(summary, courseMarker);
		if index != -1 {
			break;
		}
	}

	return summary[:index];
}