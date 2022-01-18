package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func KeepRolesUpdatedWithCourses(discord *discordgo.Session) {
	roles, err := discord.GuildRoles(GuildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return;
	}

	for courseCode, course := range courses {
		courseExistsAsRole := false;
		for _, role := range roles {
			if role.Name == courseCode {
				courseExistsAsRole = true;
			}
		}

		if !courseExistsAsRole {
			newRole, err := discord.GuildRoleCreate(GuildID);
			if err != nil {
				fmt.Println("Error creating role: ", err);
			}

			discord.GuildRoleEdit(GuildID, newRole.ID, courseCode, course.CourseColor, false, 0, true);
		}
	}
}

func GetRoleID(discord *discordgo.Session, courseCode string) string {
	roles, err := discord.GuildRoles(GuildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return "";
	}

	for _, role := range roles {
		if (role.Name == courseCode) {
			return role.ID;
		}
	}

	fmt.Println("No role for this course code found! This shouldn't happen!");
	return "";
}

func removeOtherCourseCodeRolesFromUser(discord *discordgo.Session, userID string) {
	for courseCode := range courses {
		roleID := GetRoleID(discord, courseCode);
		err := discord.GuildMemberRoleRemove(GuildID, userID, roleID)
		if err != nil {
			fmt.Println("Error removing role:", courseCode, " with roleID: ", roleID, ". Error: ", err);
		}
	}
}

func UpdateRoles(discord *discordgo.Session) {
	KeepRolesUpdatedWithCourses(discord);

	// For each user with a userID in the `users` map, change their role for the current time
	for userID, user := range users {
		fmt.Println("User has",len(user.BusyTimes), "busy times!", user.BusyTimes); 
		
		for _, busyTime := range user.BusyTimes {
			now := time.Now();

			if now.After(busyTime.start) && now.Before(busyTime.end) {
				fmt.Println("Assigning", busyTime.courseCode, " to user", user.Name);
				// Remove their course code roles
				removeOtherCourseCodeRolesFromUser(discord, userID);
				
				// Add a new role
				roleID := GetRoleID(discord, busyTime.courseCode);
				err := discord.GuildMemberRoleAdd(GuildID, userID, roleID);
				if err != nil {
					fmt.Println("Error adding role to user", userID, ". Role has ID", roleID, " and course code: ", busyTime.courseCode, ", . Error: ", err);
				}
				break;
			}
		}
	}
}