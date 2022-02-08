package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/kaspar-p/busybee/src/update"
	"github.com/pkg/errors"
)

// SLASH COMMAND CODE
// func handleCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
// 	handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]
// 	if ok {
// 		handler(discord, interaction)
// 	}
// }

func GetExternalCommandHandlers() ExternalCommandHandlerMap {
	return ExternalCommandHandlerMap{
		HandleCommand:       HandleCommand,
		BotIsReady:          BotIsReady,
		BotJoinedNewGuild:   BotJoinedNewGuild,
		BotRemovedFromGuild: BotRemovedFromGuild,
	}
}

func HandleCommand(database *persist.DatabaseType) InnerHandleCommandType {
	return func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		// Setup common error handler for panic() calls within the message handlers
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from error in handler: ", r)

				err := SendSingleMessage(
					discord,
					message.ChannelID,
					fmt.Sprintf("Error encountered while trying to handle '%s'. Please try again.", message.Content),
				)
				if err != nil {
					log.Println("Error: error sending 'recovered' message: ", err)

					return
				}
			}
		}()

		for _, command := range commandHandlers {
			contentTrigger := "." + command.Trigger

			// Check if the command matches
			if !strings.HasPrefix(message.Content, contentTrigger) {
				continue
			}

			log.Println("Matched command: ")

			// Check if the command only matches - and then garbage. e.g. .whobusybusy
			if strings.Split(message.Content, " ")[0] != contentTrigger {
				log.Println("Wrong command, prefix matched tho.")

				err := SendSingleMessage(discord, message.ChannelID, "Wrong command. Did you mean`"+contentTrigger+"`?")
				if err != nil {
					log.Println("Error: error sending 'wrong command' message: ", err)

					return
				}

				return
			}

			log.Println("Executing handler for message: ", command.Trigger)

			err := command.Handler(database, discord, message)
			if err != nil {
				log.Printf("Error encountered while executing command %s. Error: %v.\n", contentTrigger, err)

				err := SendSingleMessage(discord, message.ChannelID, "error while dealing with "+contentTrigger+" \\:(")
				if err != nil {
					log.Println("Error: error sending 'wrong command' message: ", err)

					return
				}
			}
		}
	}
}

func BotIsReady(database *persist.DatabaseType) InnerBotIsReadyType {
	return func(discord *discordgo.Session, isReady *discordgo.Ready) {
		log.Println("Bot successfully connected! Press CMD + C at any time to exit.")
		log.Println("Bot is a part of", len(isReady.Guilds), "guilds!")

		// SLASH COMMAND CODE
		// clearAndRegisterCommands(discord)

		// Once everything is ready
		update.UpdateAllGuilds(database, discord)
	}
}

func BotJoinedNewGuild(database *persist.DatabaseType) InnerBotJoinedNewGuildType {
	return func(discord *discordgo.Session, event *discordgo.GuildCreate) {
		if event.Unavailable {
			return
		} else {
			log.Println("Bot has joined a new guild with guildId: ", event.Guild.ID)
			// CREATION OF A ROLE CODE

			// Creates a role - adds it to database and GuildRoleMap
			// update.CreateRole(discord, event.Guild.ID)

			// Populate that in the users map
			// entities.Users[event.Guild.ID] = make(map[string]*entities.User)
			log.Println("\t-> Not doing anything about it, though.")
		}
	}
}

func BotRemovedFromGuild(database *persist.DatabaseType) InnerBotRemovedFromGuildType {
	return func(discord *discordgo.Session, event *discordgo.GuildDelete) {
		guildId := event.Guild.ID
		roleId := database.GetRoleIdForGuild(guildId)

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
		err := discord.GuildRoleDelete(guildId, roleId)
		if err != nil {
			log.Println("Error removing busy role from guild:", guildId, "when getting removed. Error: ", err)

			panic(errors.Wrap(err, "Error deleting role from guild"+guildId+"when getting removed. "))
		}

		// Remove all data in guild from db
		err = database.RemoveAllDataInGuild(guildId)
		if err != nil {
			log.Println("Error removing all data of a guild!")

			panic(errors.Wrap(err, "Error removing all data of a guild!"))
		}
	}
}
