package course

import (
	"fmt"
	"strings"

	"github.com/kaspar-p/bee/src/lib"

	"github.com/apognu/gocal"
)

var Courses map[string] *Course

type Course struct {
	CourseCode string
	CourseColor int
}

func InitializeCourses() {
	Courses = make(map[string] *Course);
}

func AddUnknownCourses(events []gocal.Event) {
	for _, event := range events {
		courseCode := ParseCourseCode(event.Summary);

		var decidedCourse Course;
		// If the course was already in the map - use the existing one. If not, create a new one.
		if _, ok := Courses[courseCode]; !ok {
			fmt.Println("Creating new course with code: ", courseCode);
			// Create a new course
			decidedCourse = Course{
				CourseCode: courseCode,
				CourseColor: lib.ChooseRandomColor(),
			}
			// Add the unknown course to `courses` map
			Courses[courseCode] = &decidedCourse;
		}
		
	}
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