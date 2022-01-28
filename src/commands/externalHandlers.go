package commands

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/ingest"
	"github.com/kaspar-p/bee/src/update"
	"github.com/pkg/errors"
)

// SLASH COMMAND CODE
// func handleCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
// 	handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]
// 	if ok {
// 		handler(discord, interaction)
// 	}
// }

func HandleCommand(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Setup common error handler for panic() calls within the message handlers
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from error in handler: ", r)
		}
	}()

	for key, handler := range commandHandlers {
		command := "." + key

		// Check if the command matches
		if !strings.HasPrefix(message.Content, command) {
			continue
		}

		// Check if the command only matches - and then garbage. e.g. .whobusybusy
		if strings.Split(message.Content, " ")[0] != command {
			log.Println("Wrong command, prefix matched tho.")

			err := SendSingleMessage(discord, message.ChannelID, "Wrong command. Did you mean`"+command+"`?")
			if err != nil {
				log.Println("Error: error sending 'wrong command' message: ", err)

				return
			}

			continue
		}

		log.Println("Executing handler for message: ", key)

		// Execute the handler that matches the command
		err := handler(discord, message)
		if err != nil {
			log.Printf("Error encountered while executing command %s. Error: %v.\n", command, err)

			err := SendSingleMessage(discord, message.ChannelID, "error while dealing with "+command+" \\:(")
			if err != nil {
				log.Println("Error: error sending 'wrong command' message: ", err)

				return
			}
		}

		return
	}
}

func BotIsReady(discord *discordgo.Session, isReady *discordgo.Ready) {
	log.Println("Bot successfully connected! Press CMD + C at any time to exit.")
	log.Println("Bot is a part of", len(isReady.Guilds), "guilds!")

	guildIds := make([]string, 0)
	for _, guild := range isReady.Guilds {
		guildIds = append(guildIds, guild.ID)
	}

	// Populate the users map
	entities.InitializeUsers(guildIds)

	// Connect to the database
	database.InitializeDatabase()
	ingest.FillMapsWithDatabaseData(guildIds)

	// Once everything is ready
	update.UpdateAllGuilds(discord)

	// SLASH COMMAND CODE
	// clearAndRegisterCommands(discord)

	constants.BotReady = true
}

func BotJoinedNewGuild(discord *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Unavailable || !constants.BotReady {
		return
	} else {
		log.Println("Bot has joined a new guild with guildId: ", event.Guild.ID)
		// CREATION OF A ROLE CODE

		// Creates a role - adds it to database and GuildRoleMap
		// update.CreateRole(discord, event.Guild.ID)

		// Populate that in the users map
		// entities.Users[event.Guild.ID] = make(map[string]*entities.User)
		log.Println("-> Not doing anything about it, though.")
	}
}

func BotRemovedFromGuild(discord *discordgo.Session, event *discordgo.GuildDelete) {
	guildId := event.Guild.ID
	if event.Unavailable {
		log.Println("Server has become unavailable. Guild Id: ", guildId)

		return
	}

	log.Println("Bot has been removed from guild with guildId: ", guildId)

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
	err := discord.GuildRoleDelete(guildId, update.GuildRoleMap[guildId])
	if err != nil {
		log.Println("Error removing busy role from guild:", guildId, "when getting removed. Error: ", err)

		panic(errors.Wrap(err, "Error deleting role from guild"+guildId+"when getting removed. "))
	}

	// Remove all data in guild from db
	err = database.DatabaseInstance.RemoveAllDataInGuild(guildId)
	if err != nil {
		log.Println("Error removing all data of a guild!")

		panic(errors.Wrap(err, "Error removing all data of a guild!"))
	}

	// Remove all data in memory

	// Delete users
	delete(entities.Users, guildId)

	// Delete GuildRolePair
	delete(update.GuildRoleMap, guildId)
}
