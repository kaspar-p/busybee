package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
)

func HandleWhoBusy(discord *discordgo.Session, message *discordgo.MessageCreate) {
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