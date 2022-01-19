package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/ingest"

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
	fmt.Println("Seen message: ", message.Content);
	for key, handler := range commandHandlers {
		if strings.HasPrefix(message.Content, "." + key) {
			fmt.Println("Executing handler for message: ", key);
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
		discord.ChannelMessageSend(message.ChannelID, "Requires exactly 1 `.ics` file to be attached!");
		return;
	}

	// Validate that it is a .ics file
	file := message.Attachments[0];
	if !strings.HasSuffix(file.Filename, ".ics") {
		discord.ChannelMessageSend(message.ChannelID, "Requires the file to be in `.ics` format!");
	}

	// Download the .ics file
	filepath, err := downloadFile(file.URL);
	if err != nil {
		fmt.Println("Error encountered when downloading: ", err);
		return;
	}

	// Parse the .ics file into its events
	events, err := parseCalendar(filepath);
	if errorMessage, ok := validateCalendarFile(events, err); !ok {
		discord.ChannelMessageSend(message.ChannelID, errorMessage);
		return;
	}

	fmt.Println("Going to ingest", len(events), "events!");

	// Create new users and their events
	ingest.IngestNewData(message, events);

	// Finally, update the roles when a new user is added
	UpdateRoles(discord, message.GuildID);

	discord.ChannelMessageSend(message.ChannelID, "you're enrolled \\:)")
}

func validateCalendarFile(events []gocal.Event, err error) (message string, ok bool) {
	if err != nil || len(events) == 0 {
		fmt.Println("Error encountered while parsing .ics file: ", err);
		return "Error parsing corrupt `.ics` file. No events were found \\:(", false;
	}

	// Check that all events have nonempty titles
	for _, event := range events {
		if len(event.Summary) == 0 {
			fmt.Println("Error encountered while parsing .ics file. Empty titles on some events");
			return "Error encountered while parsing .ics file. Some event's have empty titles. This is not (!) allowed \\:(", false;
		} else {
			fmt.Println("event had summary:", event.Summary);
		}
	}

	return "", true
}

func handleWhoBusy(discord *discordgo.Session, message *discordgo.MessageCreate) {
	busyUsers := make(map[string]string);
	for _, user := range entities.Users[message.GuildID] {
		if user.CurrentlyBusy.IsBusy {
			busyUsers[user.Name] = user.CurrentlyBusy.BusyWith;
		}
	}

	resultString := ""
	for name, title := range busyUsers {
		resultString = resultString + name + " is mad busy with " + title + ".\n";
	}
	if resultString == "" {
		discord.ChannelMessageSend(message.ChannelID, "no one busy \\:)");
	} else {
		discord.ChannelMessageSend(message.ChannelID, resultString);
	}
}

func handleFree(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if len(strings.Split(message.Content, " ")) != 2 {
		fmt.Println("Free command had false arguments");
		discord.ChannelMessageSend(message.ChannelID, "command must have a single argument of the @ of a user \\:)");
		return;
	}

	if len(message.Mentions) != 1 {
		fmt.Println("Free command had false arguments");
		discord.ChannelMessageSend(message.ChannelID, "command must have a single argument of the @ of a user \\:)");
		return;
	}

	mentionedID := message.Mentions[0].ID;
	
	if _, ok := entities.Users[message.GuildID][mentionedID]; !ok {
		fmt.Println("Unknown user");
		discord.ChannelMessageSend(message.ChannelID, "that user does not exist within the system. please ask them to enrol \\:)");
		return;
	}
	mentionedUser := entities.Users[message.GuildID][mentionedID];

	busyTimesToday := getTodaysEvents(mentionedUser);

	if len(busyTimesToday) == 0 {
		discord.ChannelMessageSend(message.ChannelID, "nothing going on today :)");
		return;
	}

	resultString := mentionedUser.Name + ":\n"
	for _, busyTime := range busyTimesToday {
		resultString += "    " + busyTime.Title + ": " + toNiceTimeString(busyTime.Start) + " - " + toNiceTimeString(busyTime.End) + "\n"
	}
	discord.ChannelMessageSend(message.ChannelID, resultString);
}

func toNiceTimeString(eventTime time.Time) string {
	return eventTime.Format("3:04 PM")
}

func getTodaysEvents(user *entities.User) []*entities.BusyTime {
	todaysEvents := make([]*entities.BusyTime, 0);
	for _, busyTime := range user.BusyTimes {
		now := time.Now()
		today := now.Day()
		if 	(busyTime.Start.UTC().Day() == today) && // start is today AND
			(busyTime.End.UTC().Day() == today) && // end is today AND
			(busyTime.End.After(now)) {
				todaysEvents = append(todaysEvents, busyTime);
		}
	}

	// Sort them by start time
	sort.Slice(todaysEvents, func(i, j int) bool {
		return todaysEvents[i].Start.Unix() < todaysEvents[j].Start.Unix()
	})

	return todaysEvents;
}

func botIsReady(discord *discordgo.Session, isReady *discordgo.Ready) { 
	fmt.Println("Bot successfully connected! Press CMD + C at any time to exit.");

	// Populate the ServerRoleIDMap
	for _, guild := range isReady.Guilds {
		// Get the roleID of the busy role
		KeepRolesUpdated(discord, guild.ID);

		// Populate the second level of the Users map
		entities.Users[guild.ID] = make(map[string]*entities.User);
	}
	fmt.Println("Populated role ID map from", len(isReady.Guilds), "guilds!");

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
// 			Content: "added your event information. thanks for using busybee :)",
// 		},
// 	})
// }
