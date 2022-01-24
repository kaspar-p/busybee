package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/kaspar-p/bee/src/commands"
	"github.com/kaspar-p/bee/src/constants"
	"github.com/kaspar-p/bee/src/update"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)



func init() {
	constants.InitializeViper();
	update.InitializeGuildRoleMap();
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
	discord.AddHandler(commands.BotJoinedNewGuild)
	discord.AddHandler(commands.BotRemovedFromGuild);
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildBans;

	// Open the bot
	err = discord.Open();
	if err != nil {
		fmt.Println("Error connecting to discord:", err);
		return;
	}

	// Create and start the CRON job
	cronScheduler := cron.New();
	cronScheduler.AddFunc("1 * * * * *", func() {
		update.UpdateAllGuilds(discord);
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
