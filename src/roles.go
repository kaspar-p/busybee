package main

import (
	"fmt"
	"time"

	courseLib "github.com/kaspar-p/bee/src/course"
	userLib "github.com/kaspar-p/bee/src/user"

	"github.com/bwmarrin/discordgo"
)

func KeepRolesUpdatedWithCourses(discord *discordgo.Session) {

	for courseCode, course := range courseLib.Courses {
		if _, ok := GetRoleID(discord, courseCode); !ok {
			// There was no role corresponding to the course - create one
			newRole, err := discord.GuildRoleCreate(GuildID);
			if err != nil {
				fmt.Println("Error creating role: ", err);
			}

			_, err = discord.GuildRoleEdit(GuildID, newRole.ID, courseCode, course.CourseColor, false, 0, true);
			if err != nil {
				fmt.Println("Error editing role to have the correct properties. Error: ", err);
			}
		}
	}
}

func DeleteExtraneousRoles(discord *discordgo.Session) {
	roles, err := discord.GuildRoles(GuildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return;
	}

	for courseCode := range courseLib.Courses {
		rolesWithCourseCode := []*discordgo.Role{};
		for _, role := range roles {
			if role.Name == courseCode {
				rolesWithCourseCode = append(rolesWithCourseCode, role)
			}
		}

		if len(rolesWithCourseCode) > 1 {
			// If there are MORE than 1 role with the same course code, delete all but the first one
			for _, role := range rolesWithCourseCode[1:] {
				// Remove role
				err := discord.GuildRoleDelete(GuildID, role.ID);
				if err != nil {
					fmt.Println("Error while deleting extraneous role with course code", role.Name, "and ID", role.ID, ". Error: ", err);
					break;
				}
			}
		}
	}
} 

func GetRoleID(discord *discordgo.Session, courseCode string) (string, bool) {
	roles, err := discord.GuildRoles(GuildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return "", false;
	}

	for _, role := range roles {
		if (role.Name == courseCode) {
			return role.ID, true;
		}
	}

	fmt.Println("No role for this course code found! This shouldn't happen!");
	return "", false;
}

func RemoveOtherCourseCodeRolesFromUser(discord *discordgo.Session, userID string) {
	for courseCode := range courseLib.Courses {
		roleID, ok := GetRoleID(discord, courseCode);
		if ok {
			err := discord.GuildMemberRoleRemove(GuildID, userID, roleID)
			if err != nil {
				fmt.Println("Error removing role:", courseCode, " with roleID: ", roleID, ". Error: ", err);
			}
		}
	}
}

func UpdateRoles(discord *discordgo.Session) {
	DeleteExtraneousRoles(discord);
	KeepRolesUpdatedWithCourses(discord);

	fmt.Println("Updating roles with", len(userLib.Users), "users and", len(courseLib.Courses), "courses!");

	// For each user with a userID in the `users` map, change their role for the current time
	for userID, user := range userLib.Users {
		fmt.Println("User", user.Name, "has", len(user.BusyTimes), "busy times!");

		// Remove their course code roles
		RemoveOtherCourseCodeRolesFromUser(discord, userID);
		// Assign new roles
		CheckForCurrentCourses(discord, user);
	}
}

func CheckForCurrentCourses(discord *discordgo.Session, user *userLib.User) {
	for _, busyTime := range user.BusyTimes {
		now := time.Now();
		
		// Check if now is within the bounds of the event
		if now.After(busyTime.Start) && now.Before(busyTime.End) {
			AssignCourseRoleToUser(discord, user, busyTime);
			break;
		}
	}
}

func AssignCourseRoleToUser(discord *discordgo.Session, user *userLib.User, course *userLib.BusyTime) {
	fmt.Println(user.Name, "is in", course.CourseCode, ". Assigning role to user.");

	// Add a new role
	if roleID, ok := GetRoleID(discord, course.CourseCode); ok {
		err := discord.GuildMemberRoleAdd(GuildID, user.UserID, roleID);
		if err != nil {
			fmt.Println("Error adding role to user", user.UserID, ". Role has ID", roleID, " and course code: ", course.CourseCode, ". Error: ", err);
		}
	}
}