package user

import (
	"fmt"

	"github.com/kaspar-p/bee/src/course"

	"github.com/apognu/gocal"
)

var Users map[string] *User

type User struct {
	UserID string
	Name string
	BusyTimes []*BusyTime
}

func InitializeUsers() {
	Users = make(map[string] *User);
}

func GetOrCreateUser(userID string, userName string) *User {
	if user, ok := Users[userID]; ok {
		fmt.Println("User found with ID: ", userID);
		return user;
	} else {
		fmt.Println("User created with ID:", userID);
		// Create the new user
		user := User{
			Name: userName,
			UserID: userID,
		}

		// Add the new user
		Users[userID] = &user;
		return &user;
	}
}

func (user *User) SetCourses(events []gocal.Event) {
	for _, event := range events {
		courseCode := course.ParseCourseCode(event.Summary);

		fmt.Println("Adding course " + courseCode + " to user", user.Name, ". It starts at: ", *event.Start, "and ends at", *event.End);
		busyTime := BusyTime{
			CourseCode: courseCode,
			Start: *event.Start,
			End: *event.End,
		}

		user.BusyTimes = append(user.BusyTimes, &busyTime);
	}
}