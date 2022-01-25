package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/ingest"
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
		command := "." + key;
		if strings.HasPrefix(message.Content, command)  {
			if strings.Split(message.Content, " ")[0] != command {
				fmt.Println("Wrong command, prefix matched tho.");
				discord.ChannelMessageSend(message.ChannelID, "Wrong command. Did you mean`" + command + "`?");
				return;
			}

			fmt.Println("Executing handler for message: ", key);
			err := handler(discord, message);
			if err != nil {
				fmt.Println("Error encountered while executing command:", command + ". Error: ", err);
				discord.ChannelMessageSend(message.ChannelID, "error while dealing with " + command + " \\:(");
				return;
			}
		}
	}
}

func BotIsReady(discord *discordgo.Session, isReady *discordgo.Ready) { 
	fmt.Println("Bot successfully connected! Press CMD + C at any time to exit.");
	fmt.Println("Bot is a part of", len(isReady.Guilds), "guilds!");

	guildIds := make([]string, 0);
	for _, guild := range isReady.Guilds {
		guildIds = append(guildIds, guild.ID);
	}

	// Populate the users map
	entities.InitializeUsers(guildIds);

	// Connect to the database
	database.InitializeDatabase();
	ingest.FillMapsWithDatabaseData(guildIds);

	// Once everything is ready
	update.UpdateAllGuilds(discord);

	// SLASH COMMAND CODE
	// clearAndRegisterCommands(discord);

	constants.BotReady = true;
}

func BotJoinedNewGuild(discord *discordgo.Session, event *discordgo.GuildCreate) {
	if (event.Unavailable || !constants.BotReady) {
		return;
	} else {
		fmt.Println("Bot has joined a new guild with guildId: ", event.Guild.ID);
	}

	// Creates a role - adds it to database and GuildRoleMap
	update.CreateRole(discord, event.Guild.ID);

	// Populate that in the users map
	entities.Users[event.Guild.ID] = make(map[string]*entities.User);
}

func BotRemovedFromGuild(discord *discordgo.Session, event *discordgo.GuildDelete) {
	guildId := event.Guild.ID;
	if event.Unavailable {
		fmt.Println("Server has become unavailable. Guild Id: ", guildId);
		return;
	} else {
		fmt.Println("Bot has been removed from guild with guildId: ", guildId);
	}
	
	// FROM DISCORD
	//	- remove role
	// FROM DB
	//	- remove users with guild id
	// 	- remove busy times with guild id
	//	- remove guild role pair with guild id
	// FROM MEMORY
	//	- remove users with guild id
	// 	- remove role with guild id

	// Remove role from discord
	err := discord.GuildRoleDelete(guildId, update.GuildRoleMap[guildId]);
	if err != nil {
		fmt.Println("Error removing busy role from guild:", guildId, "when getting removed. Error: ", err);
	}

	// Remove all data in guild from db
	database.DatabaseInstance.RemoveAllDataInGuild(guildId);

	// Remove all data in memory
	// TODO: rework for store structure

	// Delete users
	delete(entities.Users, guildId);

	// Delete GuildRolePair
	// TODO: rework for store[guildId] structure
	delete(update.GuildRoleMap, guildId);
}