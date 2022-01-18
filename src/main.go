package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var users map[string]*User;
var courses map[string]*Course;


var (
	GuildID string
	ChannelID string
	BotToken string
	AppID string
)

var (
	// SLASH COMMAND CODE
	// createdCommands []*discordgo.ApplicationCommand
	// commands = []*discordgo.ApplicationCommand {
	// 	{
	// 		Name: "register",
	// 		Description: "Add .ics calendar to register courses with the bot. Be sure to attach a .ics file to this message!",
	// 		Type: discordgo.ChatApplicationCommand,
	// 	},
		
	// }

	// commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.InteractionCreate) {
	// 	"register": handleEnrolment,
	// }

	commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.MessageCreate) {
		"enrol": handleEnrolment,
	}
)

func configureViper() {
	viper.SetConfigName("env");
	viper.AddConfigPath(".");
	viper.AutomaticEnv();
	viper.SetConfigType("yml");

	err := viper.ReadInConfig();
	if err != nil {
		fmt.Println("Error reading from environment variables file: ", err);
	}

	// Get environment variables
	BotToken = viper.GetString("BOT.TOKEN");
	activeServer := viper.GetString("BOT.ACTIVE_SERVER");
	GuildID = viper.GetString("BOT.GUILD_IDS." + activeServer);
	ChannelID = viper.GetString("BOT.CHANNEL_IDS." + activeServer);
	AppID = viper.GetString("BOT.APP_ID");
}

func init() {
	users = make(map[string] *User);
	courses = make(map[string] *Course);

	configureViper();
}

func main() {
	// Initialize the bot
	discord, err := discordgo.New("Bot " + BotToken);
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

	// Create and start the CRON job
	c := cron.New();
	c.AddFunc("0 * * * * *", func() {
		fmt.Println("Updating roles!");
		UpdateRoles(discord);
	});
	c.Start();
	
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

