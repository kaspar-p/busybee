package main

import (
	"fmt"
	"time"

	"github.com/kaspar-p/bee/src/constants"
	courseLib "github.com/kaspar-p/bee/src/course"
	usersLib "github.com/kaspar-p/bee/src/users"

	"github.com/bwmarrin/discordgo"
)

func KeepRolesUpdatedWithCourses(discord *discordgo.Session) {
	for courseCode, course := range courseLib.Courses {
		if _, ok := GetRoleID(discord, courseCode); !ok {
			fmt.Println("Found no role with course code", courseCode, ". Creating one.");
			// There was no role corresponding to the course - create one
			newRole, err := discord.GuildRoleCreate(constants.GuildID);
			if err != nil {
				fmt.Println("Error creating role: ", err);
			}

			_, err = discord.GuildRoleEdit(constants.GuildID, newRole.ID, courseCode, int(course.CourseColor), false, 0, true);
			if err != nil {
				fmt.Println("Error editing role to have the correct properties. Error: ", err);
			}
		}
	}
}

func DeleteExtraneousRoles(discord *discordgo.Session) {
	roles, err := discord.GuildRoles(constants.GuildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return;
	}

	for courseCode := range courseLib.Courses {
		rolesWithCourseCode := make([]*discordgo.Role, 0);
		for _, role := range roles {
			if role.Name == courseCode {
				rolesWithCourseCode = append(rolesWithCourseCode, role)
			}
		}

		if len(rolesWithCourseCode) > 1 {
			fmt.Println("Found more than one role with course code", courseCode + ". Deleting all but one.");
			// If there are MORE than 1 role with the same course code, delete all but the first one
			for _, role := range rolesWithCourseCode[1:] {
				// Remove role
				err := discord.GuildRoleDelete(constants.GuildID, role.ID);
				if err != nil {
					fmt.Println("Error while deleting extraneous role with course code", role.Name, "and ID", role.ID, ". Error: ", err);
					break;
				}
			}
		}
	}
} 

func GetRoleID(discord *discordgo.Session, courseCode string) (string, bool) {
	roles, err := discord.GuildRoles(constants.GuildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return "", false;
	}

	for _, role := range roles {
		if (role.Name == courseCode) {
			return role.ID, true;
		}
	}

	return "", false;
}

func RemoveOtherCourseCodeRolesFromUser(discord *discordgo.Session, userID string) {
	for courseCode := range courseLib.Courses {
		roleID, ok := GetRoleID(discord, courseCode);
		if ok {
			err := discord.GuildMemberRoleRemove(constants.GuildID, userID, roleID)
			if err != nil {
				fmt.Println("Error removing role:", courseCode, " with roleID: ", roleID, ". Error: ", err);
			}
		}
	}
}

func CheckForCurrentCourses(discord *discordgo.Session, user *usersLib.User) {
	for _, busyTime := range user.BusyTimes {
		now := time.Now();
		
		// Check if now is within the bounds of the event
		if now.After(busyTime.Start) && now.Before(busyTime.End) {
			user.CurrentlyBusy.IsBusy = true;
			user.CurrentlyBusy.BusyWith = busyTime.CourseCode;
			AssignCourseRoleToUser(discord, user, busyTime);
			break;
		}
	}
}

func AssignCourseRoleToUser(discord *discordgo.Session, user *usersLib.User, course *usersLib.BusyTime) {
	fmt.Println(user.Name, "is in", course.CourseCode, ". Assigning role to user.");

	// Add a new role
	if roleID, ok := GetRoleID(discord, course.CourseCode); ok {
		err := discord.GuildMemberRoleAdd(constants.GuildID, user.UserID, roleID);
		if err != nil {
			fmt.Println("Error adding role to user", user.UserID, ". Role has ID", roleID, " and course code: ", course.CourseCode, ". Error: ", err);
		}
	}
}

func UpdateRoles(discord *discordgo.Session) {
	DeleteExtraneousRoles(discord);
	KeepRolesUpdatedWithCourses(discord);

	fmt.Println("Updating roles with", len(usersLib.Users), "users and", len(courseLib.Courses), "courses!");

	// For each user with a userID in the `users` map, change their role for the current time
	for userID, user := range usersLib.Users {
		fmt.Println("User", user.Name, "has", len(user.BusyTimes), "busy times!");

		// Remove their course code roles
		RemoveOtherCourseCodeRolesFromUser(discord, userID);
		
		// Assign new roles
		CheckForCurrentCourses(discord, user);
	}
}