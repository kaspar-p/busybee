package main

import (
	"log"
	"os"
	"os/signal"

	discordLib "github.com/kaspar-p/busybee/src/discord"
	"github.com/kaspar-p/busybee/src/environment"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/kaspar-p/busybee/src/update"
	"github.com/robfig/cron"
)

func main() {
	// Initialize constants
	log.Println("Initializing constants and globals.")

	config := environment.InitializeViper(environment.PRODUCTION)

	update.InitializeGuildRoleMap()

	// Connect to the database
	database, closeDatabase := persist.InitializeDatabase(config.DatabaseConfig)
	defer closeDatabase()

	// Initialize the bot
	discord, closeDiscord := discordLib.EstablishDiscordConnection(database, config.DiscordConfig)
	defer closeDiscord()

	// Create and start the CRON job
	cronScheduler := cron.New()

	err := cronScheduler.AddFunc("1 * * * * *", func() {
		update.UpdateAllGuilds(database, discord)
	})
	if err != nil {
		log.Println("Error adding CRON job! Error: ", err)
		panic(err)
	}

	cronScheduler.Start()

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
