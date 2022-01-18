package main

import (
	"fmt"
	"strings"

	"github.com/apognu/gocal"
)

type Course struct {
	CourseCode string
	CourseColor int
}

func AddUnknownCourses(events []gocal.Event) {
	for _, event := range events {
		courseCode := ParseCourseCode(event.Summary);

		var decidedCourse Course;
		// If the course was already in the map - use the existing one. If not, create a new one.
		if _, ok := courses[courseCode]; !ok {
			fmt.Println("Creating new course with code: ", courseCode);
			// Create a new course
			decidedCourse = Course{
				CourseCode: courseCode,
				CourseColor: ChooseRandomColor(),
			}
			// Add the unknown course to `courses` map
			courses[courseCode] = &decidedCourse;
		}
		
	}
}

func ParseCourseCode(summary string) string {
	prefixes := []string { "H1", "Y1" };

	var index int;
	for _, courseMarker := range prefixes {
		index = strings.Index(summary, courseMarker);
		if index != -1 {
			break;
		}
	}

	return summary[:index];
}