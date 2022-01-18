package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	courseLib "github.com/kaspar-p/bee/src/course"
	userLib "github.com/kaspar-p/bee/src/user"

	"github.com/apognu/gocal"
	"github.com/bwmarrin/discordgo"
)

// SLASH COMMAND CODE
// func handleCommand(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
// 	handler, ok := commandHandlers[interaction.ApplicationCommandData().Name];
// 	if ok {
// 		handler(discord, interaction);
// 	}
// }

func handleCommand(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.ChannelID != ChannelID {
		return;
	}

	for key, handler := range commandHandlers {
		if strings.HasPrefix(message.Content, "." + key) {
			handler(discord, message);
		}
	}
}

func createRandomString() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	alphabet := "abcdefghijklmnopqrstuvwxyz";
	length := 20

	randBytes := make([]byte, length)
	for i := range randBytes {
		randBytes[i] = alphabet[seededRand.Intn(len(alphabet))]
	}
	return string(randBytes)
}

func handleEnrolment(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if len(message.Attachments) != 1 {
		discord.ChannelMessageSend(ChannelID, "Requires exactly 1 .ics file to be attached!");
		return;
	}

	// Validate that it is a .ics file
	file := message.Attachments[0];
	if !strings.HasSuffix(file.Filename, ".ics") {
		discord.ChannelMessageSend(ChannelID, "Requires the file to be in .ics format!");
	}

	// Download the .ics file
	filepath, err := downloadFile(file.URL);
	if err != nil {
		fmt.Println("Error encountered when downloading: ", err);
		return;
	}

	// Parse the .ics file into its events
	events, err := parseCalendar(filepath);
	if err != nil {
		fmt.Println("Error encountered while parsing .ics file: ", err);
		return;
	}

	fmt.Println("Found events (" + fmt.Sprint((len(events))) + "): ");
	for _, event := range events {
		fmt.Println(event);
	}

	courseLib.AddUnknownCourses(events);
	
	// Create a user if they do not already exist
	user := userLib.GetOrCreateUser(message.Author.ID, message.Author.Username);
	user.SetCourses(events);

	// Finally, update the roles when a new user is added
	UpdateRoles(discord);
}

func botIsReady(discord *discordgo.Session, isReady *discordgo.Ready) { 
	fmt.Println("Bot successfully connected! Press CMD + C at any time to exit.");
	// SLASH COMMAND CODE
	// clearAndRegisterCommands(discord);
}

func parseCalendar(filepath string) ([]gocal.Event, error) {
	file, err := os.Open(filepath);
	if err != nil {
		fmt.Println("Error opening .ics file: ", err);
		return nil, err;
	}
	defer file.Close();

	parser := gocal.NewParser(file);
	start, end := DetermineCurrentSemester();
	parser.Start = &start;
	parser.End = &end;
	parser.SkipBounds = true;

	err = parser.Parse();
	if err != nil {
		fmt.Println("Parsing error: ", err);
		return nil, err;
	}

	return parser.Events, nil;
}

func DetermineCurrentSemester() (time.Time, time.Time) {
	semesters := map[string]string {
		"January": "Winter",
		"February": "Winter",
		"March": "Winter",
		"April": "Winter",
		"May": "Summer",
		"June": "Summer",
		"July": "Summer",
		"August": "Summer",
		"September": "Fall",
		"October": "Fall",
		"November": "Fall",
		"December": "Fall",
	}

	now := time.Now();
	startEndMap := map[string]time.Time {
		"Winter": time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location()),
		"Summer": time.Date(now.Year(), time.May, 1, 0, 0, 0, 0, now.Location()),
		"Fall": time.Date(now.Year(), time.September, 1, 0, 0, 0, 0, now.Location()),
	}

	currentSemester := semesters[now.Month().String()];

	var start time.Time;
	var end time.Time;
	if (currentSemester == "Winter") {
		start = startEndMap["Winter"];
		end = startEndMap["Summer"];
	} else if (currentSemester == "Summer") {
		start = startEndMap["Summer"];
		end = startEndMap["Fall"];
	} else if (currentSemester == "Fall") {
		start = startEndMap["Fall"];
		end = time.Date(now.Year() + 1, time.January, 1, 0, 0, 0, 0, nil)
	}

	return start, end;
}

func downloadFile(URL string) (string, error) {
	response, err := http.Get(URL);
	if err != nil {
		fmt.Println("Error getting file: ", err);
		return "", err;
	}
	defer response.Body.Close();

	filepath := "tmp/" + createRandomString()
	output, err := os.Create(filepath);
	if err != nil {
		fmt.Printf("Error creating file at: %s. Error: %s\n", filepath, err)
		return "", err;
	}
	defer output.Close();

	_, err = io.Copy(output, response.Body);
	if err != nil {
		fmt.Println("Copying response body to file buffer failed. Error: ", err);
		return "", err;
	}

	return filepath, err;
}

// SLASH COMMAND CODE
// func clearAndRegisterCommands(discord *discordgo.Session) {
// 	// Get all global commands
// 	allGlobalCommands, err := discord.ApplicationCommands(AppID, "");
// 	if err != nil {
// 		fmt.Println("Error getting global commands: ", err);
// 	}
	
// 	// Delete all global commands associated with the bot (same ApplicationID)
// 	for _, command := range allGlobalCommands {
// 		if (AppID == command.ID) {
// 			discord.ApplicationCommandDelete(AppID, "", command.ID);
// 		}
// 	}

// 	// Get all commands in the server
// 	allCommands, err := discord.ApplicationCommands(AppID, GuildID);
// 	if err != nil {
// 		fmt.Println("Error getting slash commands: ", err);
// 		return;
// 	}

// 	// Delete all commands associated with the bot (same ApplicationID)
// 	for _, command := range allCommands {
// 		// if (AppID == command.ApplicationID) {
// 		err = discord.ApplicationCommandDelete(AppID, GuildID, command.ID);
// 		if err != nil {
// 			fmt.Println("Error deleting slash command: ", err);
// 		}
// 		// }
// 	}

// 	// Register commands again as new
// 	createdCommands, _ = discord.ApplicationCommandBulkOverwrite(discord.State.User.ID, GuildID, commands);
// 	fmt.Println("Successfully registered commands!");
// }

// func handleEnrolment(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
// 	println("Begin handling!");
// 	fmt.Println(interaction.Message);
// 	fmt.Println(interaction.Data)
// 	fmt.Println(interaction.Message.Attachments);

// 	discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: "Successfully added your course information. Thanks for using busy bee!",
// 		},
// 	})
// }
