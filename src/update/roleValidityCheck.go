package update

import (
	"log"

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
	log.Println("Processing single role ID", roleId)
	// Check if the ID is in the GuildRoleMap - if not, reset it
	if roleId != GuildRoleMap[guildId] {
		log.Println("Role ID wasn't in the GuildRoleMap, updating map.")
		ChangeGuildRoleMapEntry(guildId, roleId)
	}

	// If the role ID is NOT in discord - then create it and change the GuildRoleMap
	if !IsRoleIdInGuild(discord, guildId, roleId) {
		log.Printf("\t\t->Role ID '%s' not in discord. Creating it!\n", roleId)

		newRoleId := CreateRoleInDiscord(discord, guildId)

		ChangeGuildRoleMapEntry(guildId, newRoleId)
		database.DatabaseInstance.UpdateGuildRolePairWithNewRole(guildId, roleId, newRoleId)
	}

	log.Println("Finished processing single role ID")
}

func HandleTwoPlusRoleIds(discord *discordgo.Session, guildId string, roleIds []string) {
	chosenId := roleIds[0]
	rest := utils.RemoveStringFromSlice(roleIds, chosenId)

	// Do all of the same checks with the chosenId as if it were the only entry
	HandleOneRoleId(discord, guildId, chosenId)

	// Remove the information about the other entries
	for _, roleId := range rest {
		if IsRoleIdInGuild(discord, guildId, roleId) {
			log.Printf("Deleting extra role found in discord: %s.\n", roleId)
			DeleteRoleFromDiscord(discord, guildId, roleId)
			log.Printf("Reassigning users from role %s to role %s.\n", roleId, chosenId)
			ReassignRoles(discord, guildId, roleId, chosenId)
		}

		log.Println("Deleting extra role found in database:", roleId)
		DeleteGuildRolePairFromDatabase(guildId, roleId)
	}
}

func RunRoleValidityCheck(discord *discordgo.Session, guildId string) {
	log.Printf("\tBeginning cleanup process for guild %s!\n", guildId)

	dbRoleIds := database.DatabaseInstance.GetRoleIdsForGuild(guildId)
	log.Printf("Cleanup process found %d role IDs for this guild %s\n", len(dbRoleIds), guildId)

	switch len(dbRoleIds) {
	case 0:
		HandleZeroRoleId(discord, guildId)
	case 1:
		HandleOneRoleId(discord, guildId, dbRoleIds[0])
	default:
		HandleTwoPlusRoleIds(discord, guildId, dbRoleIds)
	}

	log.Println("Ending cleanup process!")
}
