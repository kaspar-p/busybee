package commands

import (
	"math"
	"sort"
	"strings"

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

func HandleWhoBusy(discord *discordgo.Session, message *discordgo.MessageCreate) error {
	busyUsers := make(map[string]string);
	for _, user := range entities.Users[message.GuildID] {
		if user.CurrentlyBusy.IsBusy {
			busyUsers[user.Name] = user.CurrentlyBusy.BusyWith;
		}
	}

	keys := sortKeysOfMap(busyUsers);

	embed := GenerateWhoBusyEmbed(busyUsers, keys);
	
	var err error;
	if len(busyUsers) == 0 {
		_, err = discord.ChannelMessageSendEmbed(message.ChannelID, embed);
	} else {
		_, err = discord.ChannelMessageSendEmbed(message.ChannelID, embed);
	}

	return err;
}

func GenerateWhoBusyEmbed(busyUsers map[string]string, keys []string) *discordgo.MessageEmbed {
	var lengthOfLongestName int;
	for _, name := range keys {
		lengthOfLongestName = int(math.Max(float64(len(name)), float64(lengthOfLongestName)));
	}

	resultString := "```"
	for _, name := range keys {
		spacing := strings.Repeat(" ", lengthOfLongestName - len(name));
		resultString += name + spacing + " is mad busy with " + busyUsers[name] + ".\n";
	}
	resultString += "```";

	var title string;
	var description string;
	if len(busyUsers) == 0 {
		title = "no one busy \\:)";
		description = "";
	} else {
		title = "who busy";
		description = resultString;
	}

	embed := CreateTableEmbed(title, description);
	return embed;
}