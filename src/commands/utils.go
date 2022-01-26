package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/constants"
)

func CreateTableEmbed(title, description string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Type:        "rich",
		Title:       title,
		Description: description,
		Color:       constants.BeeColor,
	}

	return &embed
}
