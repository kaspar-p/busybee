package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/update"
)

// SLASH COMMAND CODE
// func handleCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
// 	handler, ok := commandHandlers[interaction.ApplicationCommandData().Name];
// 	if ok {
// 		handler(discord, interaction);
// 	}
// }


func HandleCommand(discord *discordgo.Session, message *discordgo.MessageCreate) {
	for key, handler := range commandHandlers {
		if strings.HasPrefix(message.Content, "." + key) {
			fmt.Println("Executing handler for message: ", key);
			handler(discord, message);
		}
	}
}

func BotIsReady(discord *discordgo.Session, isReady *discordgo.Ready) { 
	fmt.Println("Bot successfully connected! Press CMD + C at any time to exit.");

	// Populate the ServerRoleIDMap
	for _, guild := range isReady.Guilds {
		// Get the roleID of the busy role
		update.KeepRolesUpdated(discord, guild.ID);

		// Populate the second level of the Users map
		entities.Users[guild.ID] = make(map[string]*entities.User);
	}
	fmt.Println("Populated role ID map from", len(isReady.Guilds), "guilds!");

	// SLASH COMMAND CODE
	// clearAndRegisterCommands(discord);
}