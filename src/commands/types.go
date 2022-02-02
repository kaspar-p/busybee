package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/persist"
)

type (
	HandleCommandType              = func(*persist.DatabaseType) InnerHandleCommandType
	InnerHandleCommandType         = func(*discordgo.Session, *discordgo.MessageCreate)
	BotIsReadyType                 = func(*persist.DatabaseType) InnerBotIsReadyType
	InnerBotIsReadyType            = func(*discordgo.Session, *discordgo.Ready)
	BotJoinedNewGuildType          = func(*persist.DatabaseType) InnerBotJoinedNewGuildType
	InnerBotJoinedNewGuildType     = func(*discordgo.Session, *discordgo.GuildCreate)
	BotRemovedFromGuildType        = func(*persist.DatabaseType) InnerBotRemovedFromGuildType
	InnerBotRemovedFromGuildType   = func(*discordgo.Session, *discordgo.GuildDelete)
	DatabaseTouchingCommandHandler = func(
		database *persist.DatabaseType,
		s *discordgo.Session,
		m *discordgo.MessageCreate) error
	PureCommandHandler = func(s *discordgo.Session, m *discordgo.MessageCreate) error
)

type ExternalCommandHandlerMap struct {
	HandleCommand       HandleCommandType
	BotIsReady          BotIsReadyType
	BotJoinedNewGuild   BotJoinedNewGuildType
	BotRemovedFromGuild BotRemovedFromGuildType
}

type CommandHandlerUnion struct {
	handlerType int
	handler     interface{}
}
