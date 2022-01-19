package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/ingest"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var (
	// SLASH COMMAND CODE
	// createdCommands []*discordgo.ApplicationCommand
	// commands = []*discordgo.ApplicationCommand {
	// 	{
	// 		Name: "register",
	// 		Description: "Add .ics calendar to register events with the bot. Be sure to attach a .ics file to this message!",
	// 		Type: discordgo.ChatApplicationCommand,
	// 	},
		
	// }

	// commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.InteractionCreate) {
	// 	"register": handleEnrolment,
	// }

	commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.MessageCreate) {
		"enrol": handleEnrolment,
		"whobusy": handleWhoBusy,
		"wyd": handleFree,
	}
)

func init() {
	constants.InitializeViper();
	entities.InitializeUsers();
	ServerRoleIDMap = make(map[string]string);
}

func main() {
	// Initialize the bot
	discord, err := discordgo.New("Bot " + constants.BotToken);
	if (err != nil) {
		fmt.Println("Error creating discord session: ", err);
		return;
	}

	// Add handlers
	discord.AddHandler(botIsReady);
	discord.AddHandler(handleCommand);
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
	defer cancel();
	
	// Create and start the CRON job
	cronScheduler := cron.New();
	cronScheduler.AddFunc("1 * * * * *", func() {
		fmt.Println("\nUpdating roles!");
		for guildID := range ServerRoleIDMap {
			fmt.Println("Updating roles for guild with id:", guildID);
			UpdateRoles(discord, guildID);
		}
		fmt.Println("Done updating roles!");
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

