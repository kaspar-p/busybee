package database

import (
	"fmt"

	coursesLib "github.com/kaspar-p/bee/src/course"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func (database *Database) AddCourse(newCourse *coursesLib.Course) string {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	result, err := database.courses.InsertOne(database.context, newCourse)
	if err != nil {
		fmt.Println("Error inserting course: ", newCourse, ". Error: ", err)
		panic(&AddCourseError{Err: err})
	}

	id := ObjectIDToString(result.InsertedID)

	return id
}

func (database *Database) GetCourses() []*coursesLib.Course {
	cursor, err := database.courses.Find(database.context, bson.D{{ }} );
	if err != nil {
		fmt.Println("Error getting cursor when finding all courses. Error: ", err);
		panic(&GetCourseError{ Err: err });
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		fmt.Println("Error getting results from cursor when getting all courses. Error: ", err);
		panic(&GetCourseError{ Err: err });
	}
	
	// Create courses out of the results
	courses := make([]*coursesLib.Course, 0);
	for _, result := range results {
		courseCode := result["coursecode"].(string)
		courseColor := int(result["coursecolor"].(int32))
		course := coursesLib.CreateCourse(courseCode, courseColor)

		courses = append(courses, course);
	}
	fmt.Println("Gotten", len(courses), "courses from database!");

	return courses;
}

func (database *Database) RemoveCourse(filter primitive.D) error {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	result, err := database.courses.DeleteOne(database.context, filter);
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %v courses.\n", result.DeletedCount);
	return nil;
}

func (database *Database) RemoveCourseByID(courseID string) {
	filter := bson.D{{Key: "_id", Value: courseID}};
	err := database.RemoveCourse(filter);
	if err != nil {
		fmt.Println("Error removing course by ID: " + courseID + ". Error: ", err);
		panic(&RemoveCourseError{Err: err});
	}
}

func (database *Database) RemoveCourseByCourseCode(courseCode string)  {
	filter := bson.D{ {Key: "CourseCode", Value: courseCode} };
	err := database.RemoveCourse(filter);
	if err != nil {
		fmt.Println("Error removing course by course code: " + courseCode + ". Error: ", err);
		panic(&RemoveCourseError{Err: err});
	}
}