package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/commands"
	"github.com/kaspar-p/bee/src/database"
	"github.com/pkg/errors"
)

type DiscordConfig struct {
	BotToken string
	AppId    string
}

func EstablishDiscordConnection(
	db *database.Database,
	config *DiscordConfig,
) (
	discord *discordgo.Session,
	disconnect func(),
) {
	// Create the bot
	discord, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		log.Printf("Error encountered while creating a bot with token %s. Err: %v\n", config.BotToken, err)

		panic(errors.Wrap(err, "Error encountered while creating a bot with token: "+config.BotToken))
	}

	// Add all of the handlers to discord
	externalHandlers := commands.GetExternalCommandHandlers()
	discord.AddHandler(externalHandlers.BotIsReady)
	discord.AddHandler(externalHandlers.HandleCommand)
	discord.AddHandler(externalHandlers.BotJoinedNewGuild)
	discord.AddHandler(externalHandlers.BotRemovedFromGuild)

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildBans

	// Open the bot
	err = discord.Open()
	if err != nil {
		log.Println("Error connecting to discord:", err)
		panic(err)
	}

	closeFunction := func() {
		err := discord.Close()
		if err != nil {
			log.Println("Error while disconnecting from discord: ", err)
			panic(errors.Wrap(err, "Error while disconnecting from discord!"))
		}
	}

	return discord, closeFunction
}
