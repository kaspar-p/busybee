package ingest

import (
	"fmt"

	courseLib "github.com/kaspar-p/bee/src/course"
	dbLib "github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/lib"
	usersLib "github.com/kaspar-p/bee/src/users"

	"github.com/apognu/gocal"
	"github.com/bwmarrin/discordgo"
)

func FillMapsWithDatabaseData() {
	// Get data and fill the `users` map
	users := dbLib.DatabaseInstance.GetUsers();
	for _, user := range users {
		usersLib.Users[user.UserID] = user;
	}

	// Get data and fill the busyTimes of each user in the `users` map
	busyTimesArray := dbLib.DatabaseInstance.GetBusyTimes();
	for _, busyTime := range busyTimesArray {
		user := usersLib.Users[busyTime.OwnerID];
		user.BusyTimes = append(user.BusyTimes, busyTime);
	}

	// Get data and fill the `courses` map
	courses := dbLib.DatabaseInstance.GetCourses();
	for _, course := range courses {
		courseLib.Courses[course.CourseCode] = course;
	}
}

func IngestNewData(message *discordgo.MessageCreate, events []gocal.Event) {
	// Update Courses from events
	CreateCoursesFromEvents(events);

	// Create a user if they do not already exist - overwrites BusyTimes
	CreateUserFromMessage(message, events);
}

func CreateUserFromMessage(message *discordgo.MessageCreate, events []gocal.Event) {
	user := GetOrCreateUser(message.Author.ID, message.Author.Username);
	
	SetCoursesFromEvents(user, events);
}

func SetCoursesFromEvents(user *usersLib.User, events []gocal.Event) {
	// Overwrite the busyTimes in memory
	for _, event := range events {
		courseCode := courseLib.ParseCourseCode(event.Summary);

		fmt.Println("Adding course " + courseCode + " to user", user.Name, ". It starts at: ", *event.Start, "and ends at", *event.End);
		busyTime := usersLib.CreateBusyTime(user.UserID, courseCode, *event.Start, *event.End);
		user.BusyTimes = append(user.BusyTimes, &busyTime);
	}

	// Overwrite the busyTimes in the database
	dbLib.DatabaseInstance.OverwriteUserBusyTimes(user.UserID, user.BusyTimes);
}

func GetOrCreateUser(userID string, userName string) *usersLib.User {
	if user, ok := usersLib.Users[userID]; ok {
		fmt.Println("User found with ID: ", userID);
		return user;
	} else {
		fmt.Println("User created with ID:", userID);
		// Create the new user
		user := usersLib.CreateUser(userName, userID);

		// Add the new user to the `users` map
		usersLib.Users[userID] = user;

		// Add the new user to the database
		dbLib.DatabaseInstance.AddUser(user);

		return user;
	}
}

func CreateCoursesFromEvents(events []gocal.Event) {
	for _, event := range events {
		courseCode := courseLib.ParseCourseCode(event.Summary);

		var course courseLib.Course;
		// If the course was already in the map - use the existing one. If not, create a new one.
		if _, ok := courseLib.Courses[courseCode]; !ok {
			fmt.Println("Creating new course with code: ", courseCode);
			// Create a new course
			course = courseLib.Course{
				CourseCode: courseCode,
				CourseColor: lib.ChooseRandomColor(),
			}

			// Add the unknown course to `courses` map
			courseLib.Courses[courseCode] = &course;

			// Add the unknown course to the database
			dbLib.DatabaseInstance.AddCourse(&course);
		}
	}
}