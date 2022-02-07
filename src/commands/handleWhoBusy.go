package commands

import (
	"math"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/kaspar-p/busybee/src/utils"
	"github.com/pkg/errors"
)

func sortKeysOfMap(unsortedMap map[string]string) []string {
	// Get keys of map
	keys := make([]string, 0)
	for key := range unsortedMap {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return unsortedMap[keys[i]] < unsortedMap[keys[j]]
	})

	return keys
}

func HandleWhoBusy(database *persist.DatabaseType, discord *discordgo.Session, message *discordgo.MessageCreate) error {
	busyUsers := database.GetCurrentlyBusyUsers()

	keys := sortKeysOfMap(busyUsers)
	embed := GenerateWhoBusyEmbed(busyUsers, keys)

	_, err := discord.ChannelMessageSendEmbed(message.ChannelID, embed)

	return errors.Wrap(err, "Error sending response to .whobusy message.")
}

func GenerateWhoBusyEmbed(busyUsers map[string]string, keys []string) *discordgo.MessageEmbed {
	var lengthOfLongestName int
	for _, name := range keys {
		lengthOfLongestName = int(math.Max(float64(len(name)), float64(lengthOfLongestName)))
	}

	resultString := ""

	for _, name := range keys {
		spacing := strings.Repeat(" ", lengthOfLongestName-len(name))
		resultString += name + spacing + " is mad busy with " + busyUsers[name] + ".\n"
	}

	var title, description string

	if len(busyUsers) == 0 {
		title = "no one busy \\:)"
		description = ""
	} else {
		title = "who busy"
		description = utils.WrapStringInCodeBlock(resultString)
	}

	embed := CreateGenericEmbed(title, description)

	return embed
}
