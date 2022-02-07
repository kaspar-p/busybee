package commands

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apognu/gocal"
	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/ingest"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/kaspar-p/busybee/src/update"
	"github.com/pkg/errors"
)

func HandleEnrol(database *persist.DatabaseType, discord *discordgo.Session, message *discordgo.MessageCreate) error {
	if len(message.Attachments) != 1 {
		err := SendSingleMessage(discord, message.ChannelID, "Requires exactly 1 `.ics` file to be attached!")

		return err
	}

	// Validate that it is a .ics file
	file := message.Attachments[0]
	if !strings.HasSuffix(file.Filename, ".ics") {
		err := SendSingleMessage(discord, message.ChannelID, "Requires the file to be in `.ics` format!")

		return err
	}

	// Download the .ics file
	filepath, err := downloadFile(file.URL)
	if err != nil {
		log.Println("Error encountered when downloading: ", err)

		sendErr := SendSingleMessage(discord, message.ChannelID, "Downloading .ics file failed. Please try again.")

		return errors.Wrap(err, sendErr.Error())
	}

	// Cleanup the file that was created
	defer removeFile(filepath)

	// Parse the .ics file into its events
	events, err := parseCalendar(filepath)
	if errorMessage, ok := validateCalendarFile(events, err); !ok {
		sendErr := SendSingleMessage(discord, message.ChannelID, errorMessage)

		return errors.Wrap(err, sendErr.Error())
	}

	log.Printf("Going to ingest %d events!\n", len(events))

	// Create new users and their events
	ingest.IngestNewData(database, message, events)

	// Finally, update the roles when a new user is added
	update.UpdateSingleGuild(discord, message.GuildID)

	err = SendSingleMessage(discord, message.ChannelID, "you're enrolled \\:)")

	return err
}

func createRandomString() string {
	stringLength := 20
	randomBytes := make([]byte, stringLength)

	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return base32.StdEncoding.EncodeToString(randomBytes)[:stringLength]
}

func validateCalendarFile(events []gocal.Event, err error) (message string, ok bool) {
	if err != nil || len(events) == 0 {
		log.Println("Error encountered while parsing .ics file: ", err)

		return "Error parsing corrupt `.ics` file. No events were found \\:(", false
	}

	// Check that all events have nonempty titles
	for i := 0; i < len(events); i++ {
		event := events[i]

		if event.Summary == "" {
			log.Println("Error encountered while parsing .ics file. Empty titles on some events")

			return "Error encountered while parsing .ics file. " +
				"Some event's have empty titles. This is not (!) allowed \\:(", false
		}
	}

	return "", true
}

func parseCalendar(filepath string) ([]gocal.Event, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Error opening .ics file: ", err)

		return nil, errors.Wrap(err, "Error opening .ics file.")
	}
	defer file.Close()

	parser := gocal.NewParser(file)
	start, end := DetermineCurrentSemester()
	parser.Start = &start
	parser.End = &end
	parser.SkipBounds = true

	err = parser.Parse()
	if err != nil {
		log.Println("Parsing error: ", err)

		return nil, errors.Wrap(err, "Parsing .ics file error.")
	}

	return parser.Events, nil
}

func DetermineCurrentSemester() (start, end time.Time) {
	semesters := map[string]string{
		"January":   "Winter",
		"February":  "Winter",
		"March":     "Winter",
		"April":     "Winter",
		"May":       "Summer",
		"June":      "Summer",
		"July":      "Summer",
		"August":    "Summer",
		"September": "Fall",
		"October":   "Fall",
		"November":  "Fall",
		"December":  "Fall",
	}

	now := time.Now()
	startEndMap := map[string]time.Time{
		"Winter": time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location()),
		"Summer": time.Date(now.Year(), time.May, 1, 0, 0, 0, 0, now.Location()),
		"Fall":   time.Date(now.Year(), time.September, 1, 0, 0, 0, 0, now.Location()),
	}

	currentSemester := semesters[now.Month().String()]

	switch currentSemester {
	case "Winter":
		start = startEndMap["Winter"]
		end = startEndMap["Summer"]
	case "Summer":
		start = startEndMap["Summer"]
		end = startEndMap["Fall"]
	case "Fall":
		start = startEndMap["Fall"]
		end = time.Date(now.Year()+1, time.January, 1, 0, 0, 0, 0, nil)
	}

	return start, end
}

func removeFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		log.Println("Error removing file with path", filepath)
		panic(errors.New("Error removing file with path: " + filepath))
	}
}

func downloadFile(url string) (string, error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Println("Error setting up request:", err)

		return "", errors.Wrap(err, "Error setting up request to url: "+url)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting file: ", err)

		return "", errors.Wrap(err, "Error getting the download URL from discord. Didn't download.")
	}
	defer response.Body.Close()

	var perms fs.FileMode = 0o0700
	err = os.Mkdir("./tmp/", perms)

	if err != nil && !os.IsExist(err) {
		log.Println("Error while mkdir tmp: ", err)

		return "", errors.Wrap(err, "Error encountered while making temp directory!")
	} else if os.IsExist(err) {
		log.Println("Skipping creating temp directory - already exists.")
	}

	filepath := "./tmp/" + createRandomString()

	output, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating file at: %s. Error: %s\n", filepath, err)

		return "", errors.Wrap(err, "Error creating file.")
	}
	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Println("Copying response body to file buffer failed. Error: ", err)

		return "", errors.Wrap(err, "Error copying to output file from output stream.")
	}

	return filepath, errors.Wrap(err, "Error downloading file!")
}

// SLASH COMMAND CODE
// func clearAndRegisterCommands(discord *discordgo.Session) {
// 	// Get all global commands
// 	allGlobalCommands, err := discord.ApplicationCommands(AppID, "")
// 	if err != nil {
// 		log.Panic("Error getting global commands: ", err)
// 	}

// 	// Delete all global commands associated with the bot (same ApplicationID)
// 	for _, command := range allGlobalCommands {
// 		if (AppID == command.ID) {
// 			discord.ApplicationCommandDelete(AppID, "", command.ID)
// 		}
// 	}

// 	// Get all commands in the guild
// 	allCommands, err := discord.ApplicationCommands(AppID, GuildID)
// 	if err != nil {
// 		log.Panic("Error getting slash commands: ", err)
// 		return
// 	}

// 	// Delete all commands associated with the bot (same ApplicationID)
// 	for _, command := range allCommands {
// 		// if (AppID == command.ApplicationID) {
// 		err = discord.ApplicationCommandDelete(AppID, GuildID, command.ID)
// 		if err != nil {
// 			log.Panic("Error deleting slash command: ", err)
// 		}
// 		// }
// 	}

// 	// Register commands again as new
// 	createdCommands, _ = discord.ApplicationCommandBulkOverwrite(discord.State.User.ID, GuildID, commands)
// 	log.Println("Successfully registered commands!")
// }

// func handleEnrolment(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
// 	println("Begin handling!")
// 	log.Println(interaction.Message)
// 	log.Println(interaction.Data)
// 	log.Println(interaction.Message.Attachments)

// 	discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: "added your event information. thanks for using busybee :)",
// 		},
// 	})
// }
