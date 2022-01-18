package main

import (
	"fmt"
	"time"

	"github.com/apognu/gocal"
)


type User struct {
	UserID string
	Name string
	BusyTimes []BusyTime
}

type BusyTime struct {
	courseCode string
	start time.Time
	end time.Time
}

func GetOrCreateUser(userID string, userName string) *User {
	if user, ok := users[userID]; ok {
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
		users[userID] = &user;
		return &user;
	}
}

func (user *User) SetCourses(events []gocal.Event) {
	for _, event := range events {
		courseCode := ParseCourseCode(event.Summary);

		fmt.Println("Adding course " + courseCode + " to user", user.Name, ". It starts at: ", *event.Start, "and ends at", *event.End);
		busyTime := BusyTime{
			courseCode: courseCode,
			start: *event.Start,
			end: *event.End,
		}

		user.BusyTimes = append(user.BusyTimes, busyTime);
	}
}