package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/pkg/errors"
)

func toNiceTimeString(eventTime time.Time) string {
	return eventTime.Format("3:04 PM")
}

func HandleWyd(discord *discordgo.Session, message *discordgo.MessageCreate) error {
	expectedArgumentNum := 2
	if len(strings.Split(message.Content, " ")) != expectedArgumentNum {
		fmt.Println("Free command had false arguments")
		discord.ChannelMessageSend(message.ChannelID,
			"command must have a single argument of the @ of a user \\:)")

		return nil
	}

	expectedMentions := 1
	if len(message.Mentions) != expectedMentions {
		fmt.Println("Free command had false arguments")
		discord.ChannelMessageSend(message.ChannelID,
			"command must have a single argument of the @ of a user \\:)")

		return nil
	}

	mentionedId := message.Mentions[0].ID

	if mentionedId == discord.State.User.ID {
		fmt.Println("Asked the bot wyd!")
		discord.ChannelMessageSend(message.ChannelID, "nothing much \\;)")

		return nil
	}

	if _, ok := entities.Users[message.GuildID][mentionedId]; !ok {
		fmt.Println("Unknown user")
		discord.ChannelMessageSend(message.ChannelID,
			"that user does not exist within the system. please ask them to enrol \\:)")

		return nil
	}

	mentionedUser := entities.Users[message.GuildID][mentionedId]
	busyTimesToday := mentionedUser.GetTodaysEvents()

	if len(busyTimesToday) == 0 {
		discord.ChannelMessageSend(message.ChannelID, "nothing going on for the rest of today :)")

		return nil
	}

	embed := GenerateWydEmbed(busyTimesToday, mentionedUser)

	_, err := discord.ChannelMessageSendEmbed(message.ChannelID, embed)

	return errors.Wrap(err, "Error sending response to .wyd message.")
}

func GenerateWydEmbed(busyTimesToday []*entities.BusyTime, mentionedUser *entities.User) *discordgo.MessageEmbed {
	resultString := "```"
	for _, busyTime := range busyTimesToday {
		resultString += busyTime.Title + ": " + toNiceTimeString(busyTime.Start) +
			" - " + toNiceTimeString(busyTime.End) + "\n"
	}

	resultString += "```"

	embed := CreateTableEmbed(mentionedUser.Name, resultString)

	return embed
}
