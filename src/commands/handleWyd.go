package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
)


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


func HandleWyd(discord *discordgo.Session, message *discordgo.MessageCreate) {
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

	if mentionedID == discord.State.User.ID {
		fmt.Println("Asked the bot wyd!")
		discord.ChannelMessageSend(message.ChannelID, "nothing much \\;)");
		return;
	}
	
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