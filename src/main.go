package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/commands"
	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/update"
	"github.com/robfig/cron"
)

func main() {
	// Initialize constants
	log.Println("Initializing constants and globals.")

	constants.InitializeViper()
	update.InitializeGuildRoleMap()

	// Initialize the bot
	discord, err := discordgo.New("Bot " + constants.BotToken)
	if err != nil {
		log.Println("Error creating discord session: ", err)
		panic(err)
	}

	// Add handlers
	discord.AddHandler(commands.BotIsReady)
	discord.AddHandler(commands.HandleCommand)
	discord.AddHandler(commands.BotJoinedNewGuild)
	discord.AddHandler(commands.BotRemovedFromGuild)
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildBans

	// Open the bot
	err = discord.Open()
	if err != nil {
		log.Println("Error connecting to discord:", err)
		panic(err)
	}

	// Create and start the CRON job
	cronScheduler := cron.New()

	err = cronScheduler.AddFunc("1 * * * * *", func() {
		update.UpdateAllGuilds(discord)
	})
	if err != nil {
		log.Println("Error adding CRON job! Error: ", err)
		panic(err)
	}

	cronScheduler.Start()

	defer discord.Close()

	// SLASH COMMAND CODE
	// // Remove the commands when the bot is closed
	// for _, command := range createdCommands {
	// 	err := discord.ApplicationCommandDelete(discord.State.User.ID, GuildId, command.ID)
	// 	if err != nil {
	// 		log.Println("Cannot delete %q command: %v", command.Name, err)
	// 	}
	// }

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
