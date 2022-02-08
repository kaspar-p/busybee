package commands

import (
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/entities"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/pkg/errors"
)

func toNiceTimeString(eventTime time.Time) string {
	return eventTime.Format("3:04 PM")
}

func validateWydStructure(discord *discordgo.Session, message *discordgo.MessageCreate) error {
	expectedArgumentNum := 2
	if len(strings.Split(message.Content, " ")) != expectedArgumentNum {
		log.Println("Free command had false arguments")

		err := SendSingleMessage(discord, message.ChannelID, "command must have a single argument of the @ of a user \\:)")

		return err
	}

	expectedMentions := 1
	if len(message.Mentions) != expectedMentions {
		log.Println("Free command had false arguments")

		err := SendSingleMessage(discord, message.ChannelID, "command must have a single argument of the @ of a user \\:)")

		return err
	}

	return nil
}

func HandleWyd(database *persist.DatabaseType, discord *discordgo.Session, message *discordgo.MessageCreate) error {
	err := validateWydStructure(discord, message)
	if err != nil {
		return err
	}

	mentionedId := message.Mentions[0].ID

	if mentionedId == discord.State.User.ID {
		return TalkToBusyBee(discord, message, "wyd")
	}

	user, userExists := database.GetUser(message.GuildID, mentionedId)
	if !userExists {
		log.Printf("Requirement failed - sending error message %s.\n", "user DNE")

		err = SendSingleMessage(discord,
			message.ChannelID,
			"that user does not exist within the system. please ask them to enrol \\:)",
		)

		return err
	}

	busyTimesToday := database.GetTodaysEventsForUser(user.Id)

	if len(busyTimesToday) == 0 {
		err = SendSingleMessage(discord, message.ChannelID, "nothing going on for the rest of today :)")

		return err
	}

	embed := GenerateWydEmbed(busyTimesToday, user)
	_, err = discord.ChannelMessageSendEmbed(message.ChannelID, embed)

	return errors.Wrap(err, "Error sending response to .wyd message.")
}

func GenerateWydEmbed(busyTimesToday []*entities.BusyTime, mentionedUser *entities.User) *discordgo.MessageEmbed {
	resultString := "```"
	for _, busyTime := range busyTimesToday {
		resultString += busyTime.Title + ": " + toNiceTimeString(busyTime.Start) +
			" - " + toNiceTimeString(busyTime.End) + "\n"
	}

	resultString += "```"

	embed := CreateGenericEmbed(mentionedUser.Name, resultString)

	return embed
}
