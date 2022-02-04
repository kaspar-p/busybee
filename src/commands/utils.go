package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/busybee/src/constants"
	"github.com/pkg/errors"
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

func TalkToBusyBee(discord *discordgo.Session, message *discordgo.MessageCreate, command string) error {
	log.Println("Someone's talking to busybee through command:", command)

	var response string

	switch command {
	case ".wyd":
		response = "nothing much, wbu \\;)"
	case ".whenfree":
		response = "i'm always free for you \\;)"
	}

	err := SendSingleMessage(
		discord,
		message.ChannelID,
		response,
	)

	return err
}

func SendSingleMessage(discord *discordgo.Session, channelID, contents string) error {
	_, err := discord.ChannelMessageSend(channelID, contents)
	if err != nil {
		errorMessage := fmt.Sprintf("Error sending single message with contents %s: %v", contents, err)
		log.Println(errorMessage)

		return errors.Wrap(err, errorMessage)
	}

	return nil
}

func SendSingleEmbed(discord *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) error {
	_, err := discord.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		errorMessage := fmt.Sprintf("Error sending single message with embed %v: %v", embed, err)
		log.Println(errorMessage)

		return errors.Wrap(err, errorMessage)
	}

	return nil
}
