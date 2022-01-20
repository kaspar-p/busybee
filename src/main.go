package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kaspar-p/bee/src/commands"
	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/ingest"
	"github.com/kaspar-p/bee/src/update"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)



func init() {
	constants.InitializeViper();
	entities.InitializeUsers();
	update.InitializeServerRoleIDMap();
	
}

func main() {
	// Initialize the bot
	discord, err := discordgo.New("Bot " + constants.BotToken);
	if (err != nil) {
		fmt.Println("Error creating discord session: ", err);
		return;
	}

	// Add handlers
	discord.AddHandler(commands.BotIsReady);
	discord.AddHandler(commands.HandleCommand);
	discord.Identify.Intents = discordgo.IntentsGuildMessages;

	// Open the bot
	err = discord.Open();
	if err != nil {
		fmt.Println("Error connecting to discord:", err);
		return;
	}

	// Connect to the database
	cancel := database.InitializeDatabase();
	ingest.FillMapsWithDatabaseData();

	// Once everything is ready
	update.UpdateAllServers(discord);
	defer cancel();
	
	// Create and start the CRON job
	cronScheduler := cron.New();
	cronScheduler.AddFunc("1 * * * * *", func() {
		update.UpdateAllServers(discord);
	});
	cronScheduler.Start();
	
	defer discord.Close();
	stop := make(chan os.Signal, 1);
	signal.Notify(stop, os.Interrupt);
	<-stop;

	// SLASH COMMAND CODE
	// // Remove the commands when the bot is closed
	// for _, command := range createdCommands {
	// 	err := discord.ApplicationCommandDelete(discord.State.User.ID, GuildID, command.ID)
	// 	if err != nil {
	// 		log.Fatalf("Cannot delete %q command: %v", command.Name, err)
	// 	}
	// }
}
