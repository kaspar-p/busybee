package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
)


func toNiceTimeString(eventTime time.Time) string {
	return eventTime.Format("3:04 PM")
}


func HandleWyd(discord *discordgo.Session, message *discordgo.MessageCreate) error {
	if len(strings.Split(message.Content, " ")) != 2 {
		fmt.Println("Free command had false arguments");
		discord.ChannelMessageSend(message.ChannelID, "command must have a single argument of the @ of a user \\:)");
		return nil;
	}

	if len(message.Mentions) != 1 {
		fmt.Println("Free command had false arguments");
		discord.ChannelMessageSend(message.ChannelID, "command must have a single argument of the @ of a user \\:)");
		return nil;
	}

	mentionedID := message.Mentions[0].ID;

	if mentionedID == discord.State.User.ID {
		fmt.Println("Asked the bot wyd!")
		discord.ChannelMessageSend(message.ChannelID, "nothing much \\;)");
		return nil;
	}
	
	if _, ok := entities.Users[message.GuildID][mentionedID]; !ok {
		fmt.Println("Unknown user");
		discord.ChannelMessageSend(message.ChannelID, "that user does not exist within the system. please ask them to enrol \\:)");
		return nil;
	}
	mentionedUser := entities.Users[message.GuildID][mentionedID];
	busyTimesToday := mentionedUser.GetTodaysEvents();

	if len(busyTimesToday) == 0 {
		discord.ChannelMessageSend(message.ChannelID, "nothing going on for the rest of today :)");
		return nil;
	}

	resultString := mentionedUser.Name + ":\n"
	for _, busyTime := range busyTimesToday {
		resultString += "    " + busyTime.Title + ": " + toNiceTimeString(busyTime.Start) + " - " + toNiceTimeString(busyTime.End) + "\n"
	}
	
	_, err := discord.ChannelMessageSend(message.ChannelID, resultString);
	return err;
}