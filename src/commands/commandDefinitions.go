package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/database"
)

// SLASH COMMAND CODE
// createdCommands []*discordgo.ApplicationCommand
// commands = []*discordgo.ApplicationCommand {
// 	{
// 		Name: "register",
// 		Description: "Add .ics calendar to register events with the bot. Be sure to attach a .ics file to this message!",
// 		Type: discordgo.ChatApplicationCommand,
// 	},.

// }.

// commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.InteractionCreate) {
// 	"register": handleEnrolment,
// }.

type (
	HandleCommandType            = func(*database.Database) InnerHandleCommandType
	InnerHandleCommandType       = func(*discordgo.Session, *discordgo.MessageCreate)
	BotIsReadyType               = func(*database.Database) InnerBotIsReadyType
	InnerBotIsReadyType          = func(*discordgo.Session, *discordgo.Ready)
	BotJoinedNewGuildType        = func(*database.Database) InnerBotJoinedNewGuildType
	InnerBotJoinedNewGuildType   = func(*discordgo.Session, *discordgo.GuildCreate)
	BotRemovedFromGuildType      = func(*database.Database) InnerBotRemovedFromGuildType
	InnerBotRemovedFromGuildType = func(*discordgo.Session, *discordgo.GuildDelete)
	CommandHandler               = func(s *discordgo.Session, m *discordgo.MessageCreate) error
)

type ExternalCommandHandlerMap struct {
	HandleCommand       HandleCommandType
	BotIsReady          BotIsReadyType
	BotJoinedNewGuild   BotJoinedNewGuildType
	BotRemovedFromGuild BotRemovedFromGuildType
}

var commandHandlers = map[string]CommandHandler{
	"enrol":    HandleEnrol,
	"whobusy":  HandleWhoBusy,
	"wyd":      HandleWyd,
	"whenfree": HandleWhenFree,
}
