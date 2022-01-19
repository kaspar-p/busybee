package main

import (
	"fmt"
	"time"

	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/entities"

	"github.com/bwmarrin/discordgo"
)

var ServerRoleIDMap map[string]string;

func KeepRolesUpdated(discord *discordgo.Session, guildID string) {
	roles, err := discord.GuildRoles(guildID);
	if err != nil {
		fmt.Println("Error getting roles: ", err);
		return;
	}

	foundBusyRoleID := "";
	for _, role := range roles {
		if role.Name == constants.BusyRoleName {
			foundBusyRoleID = role.ID;
		}
	}

	if foundBusyRoleID != "" {
		ServerRoleIDMap[guildID] = foundBusyRoleID;
	} else {
		fmt.Println("Found no busy role. Creating one.");
		
		// There was no "busy" role - create one
		newRole, err := discord.GuildRoleCreate(guildID);
		if err != nil {
			fmt.Println("Error creating role: ", err);
		}

		// Set global busy role ID to be the new ID
		ServerRoleIDMap[guildID] = newRole.ID;

		_, err = discord.GuildRoleEdit(guildID, newRole.ID, constants.BusyRoleName, 12847710, false, 0, true);
		if err != nil {
			fmt.Println("Error editing role to have the correct properties. Error:", err);
		}
	}
}

func CheckIfUserBusy(discord *discordgo.Session, user *entities.User, guildID string) {
	stayBusy := false;
	for _, busyTime := range user.BusyTimes {
		now := time.Now();
		
		// Check if now is within the bounds of the event
		if now.After(busyTime.Start) && now.Before(busyTime.End) && !user.CurrentlyBusy.IsBusy && user.CurrentlyBusy.BusyWith != busyTime.Title {
			stayBusy = true;
			user.CurrentlyBusy.IsBusy = true;
			user.CurrentlyBusy.BusyWith = busyTime.Title;
			
			err := discord.GuildMemberRoleAdd(guildID, user.ID, ServerRoleIDMap[guildID]);
			if err != nil {
				fmt.Println("Error adding role to user", user.ID, ", and title:", user.CurrentlyBusy.BusyWith, ". Error:", err);
			}
			break;
		}
	}

	if !stayBusy {
		user.CurrentlyBusy.BusyWith = "";
		user.CurrentlyBusy.IsBusy = false;
		err := discord.GuildMemberRoleRemove(guildID, user.ID, ServerRoleIDMap[guildID]);
		if err != nil {
			fmt.Println("Error removing role from user", user.ID, ". Error: ", err);
		}
	}
}

func UpdateRoles(discord *discordgo.Session, guildID string) {
	KeepRolesUpdated(discord, guildID);

	fmt.Println("Updating roles with", len(entities.Users[guildID]), "users in this server!");
	// For each user with a ID in the `users` map, change their role for the current time
	for _, user := range entities.Users[guildID] {
		if user.BelongsTo == guildID {
			fmt.Println("User", user.Name, "has", len(user.BusyTimes), "busy times!");

			// Check if the user should change busy status
			CheckIfUserBusy(discord, user, guildID);
		} else {
			fmt.Println("This shouldn't happen! User", user.Name, " has guildID", user.BelongsTo, "but should have", guildID);
		}
	}
}