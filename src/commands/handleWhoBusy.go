package commands

import (
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
)

func sortKeysOfMap(unsortedMap map[string]string) []string {
	// Get keys of map
	keys := make([]string, 0);
	for key := range unsortedMap {
		keys = append(keys, key);
	}

	sort.Slice(keys, func(i, j int) bool {
		return unsortedMap[keys[i]] < unsortedMap[keys[j]]
	})

	return keys;
}

func HandleWhoBusy(discord *discordgo.Session, message *discordgo.MessageCreate) {
	busyUsers := make(map[string]string);
	for _, user := range entities.Users[message.GuildID] {
		if user.CurrentlyBusy.IsBusy {
			busyUsers[user.Name] = user.CurrentlyBusy.BusyWith;
		}
	}

	keys := sortKeysOfMap(busyUsers);

	resultString := ""
	for _, name := range keys {
		resultString = resultString + name + " is mad busy with " + busyUsers[name] + ".\n";
	}
	if resultString == "" {
		discord.ChannelMessageSend(message.ChannelID, "no one busy \\:)");
	} else {
		discord.ChannelMessageSend(message.ChannelID, resultString);
	}
}