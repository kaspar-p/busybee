package update

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/utils"
)

var GuildRoleMap map[string]string

func InitializeGuildRoleMap() {
	GuildRoleMap = make(map[string]string)
}

func HandleZeroRoleId(discord *discordgo.Session, guildId string) {
	CreateRole(discord, guildId)
}

func HandleOneRoleId(discord *discordgo.Session, guildId, roleId string) {
	fmt.Println("Processing single role ID", roleId)
	// Check if the ID is in the GuildRoleMap - if not, reset it
	if roleId != GuildRoleMap[guildId] {
		fmt.Println("Role ID wasn't in the GuildRoleMap, updating map.")
		ChangeGuildRoleMapEntry(guildId, roleId)
	}

	// If the role ID is NOT in discord - then create it and change the GuildRoleMap
	if !IsRoleIdInGuild(discord, guildId, roleId) {
		fmt.Printf("\t\t->Role ID '%s' not in discord. Creating it!\n", roleId)

		newRoleId := CreateRoleInDiscord(discord, guildId)

		ChangeGuildRoleMapEntry(guildId, newRoleId)
		database.DatabaseInstance.UpdateGuildRolePairWithNewRole(guildId, roleId, newRoleId)
	}

	fmt.Println("Finished processing single role ID")
}

func HandleTwoPlusRoleIds(discord *discordgo.Session, guildId string, roleIds []string) {
	chosenId := roleIds[0]
	rest := utils.RemoveStringFromSlice(roleIds, chosenId)

	// Do all of the same checks with the chosenId as if it were the only entry
	HandleOneRoleId(discord, guildId, chosenId)

	// Remove the information about the other entries
	for _, roleId := range rest {
		if IsRoleIdInGuild(discord, guildId, roleId) {
			fmt.Println("Deleting extra role found in discord:", roleId)
			DeleteRoleFromDiscord(discord, guildId, roleId)
			fmt.Printf("Reassigning users from role %s to role %s.\n", roleId, chosenId)
			ReassignRoles(discord, guildId, roleId, chosenId)
		}

		fmt.Println("Deleting extra role found in database:", roleId)
		DeleteGuildRolePairFromDatabase(guildId, roleId)
	}
}

func RunRoleValidityCheck(discord *discordgo.Session, guildId string) {
	fmt.Printf("\tBeginning cleanup process for guild %s!\n", guildId)

	dbRoleIds := database.DatabaseInstance.GetRoleIdsForGuild(guildId)
	fmt.Printf("Cleanup process found %d role IDs for this guild %s\n", len(dbRoleIds), guildId)

	switch len(dbRoleIds) {
	case 0:
		HandleZeroRoleId(discord, guildId)
	case 1:
		HandleOneRoleId(discord, guildId, dbRoleIds[0])
	default:
		HandleTwoPlusRoleIds(discord, guildId, dbRoleIds)
	}

	fmt.Println("Ending cleanup process!")
}
