package commands

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

var (
	DatabaseTouchingCommandHandlerType = 1
	PureCommandHandlerType             = 2
)

func (u CommandHandlerUnion) unionToPureCommandHandler() PureCommandHandler {
	if handler, ok := u.handler.(PureCommandHandler); ok {
		return handler
	}

	return nil
}

func (u CommandHandlerUnion) unionToDatabaseTouchingCommandHandler() DatabaseTouchingCommandHandler {
	if handler, ok := u.handler.(DatabaseTouchingCommandHandler); ok {
		return handler
	}

	return nil
}

var commandHandlers = map[string]CommandHandlerUnion{
	"enrol":    {DatabaseTouchingCommandHandlerType, HandleEnrol},
	"whobusy":  {PureCommandHandlerType, HandleWhoBusy},
	"wyd":      {PureCommandHandlerType, HandleWyd},
	"whenfree": {PureCommandHandlerType, HandleWhenFree},
}
